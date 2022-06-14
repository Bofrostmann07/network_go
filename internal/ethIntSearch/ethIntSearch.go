package ethIntSearch

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"network_go/internal/config"
	"network_go/internal/models"
	"network_go/internal/util/ioUtil"
	"network_go/internal/util/sshUtil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func fetchEthIntConfig(switchInventory *[]models.NetworkSwitch) {
	CheckConnection(switchInventory)
	if switchInventory == nil {
		log.Fatalln("Couldn't reach any switch")
	}

	// Filter all unreachable switches out
	var reachableSwitchList []models.NetworkSwitch
	for _, networkSwitch := range *switchInventory {
		if networkSwitch.Reachable {
			reachableSwitchList = append(reachableSwitchList, networkSwitch)
		}
	}

	jobs := make(chan models.NetworkSwitch, len(reachableSwitchList))
	results := make(chan map[string]models.EthInterface, len(reachableSwitchList))
	command := "show config"
	for i := 0; i < config.AppConfig.SSH.MaxGoRoutines; i++ {
		go worker2(jobs, command, results)
	}
	for _, networkSwitch := range reachableSwitchList {
		jobs <- networkSwitch
	}
	close(jobs)

	for i := 0; i < len(reachableSwitchList); i++ {
		(*switchInventory)[i].EthInterfaces = <-results
		fmt.Println(<-results)
	}
	close(results)

	saveAsJson(switchInventory)
}

func worker2(jobs <-chan models.NetworkSwitch, command string, results chan<- map[string]models.EthInterface) {
	for networkSwitch := range jobs {
		rawOutput := sendSingleShowCommand(command, networkSwitch)
		results <- parseIntEthConfig(rawOutput)

	}
}

func sendSingleShowCommand(command string, networkSwitch models.NetworkSwitch) (rawOutput string) {
	addr := net.JoinHostPort(networkSwitch.Address, config.AppConfig.SSH.Port)
	rawOutput = sshUtil.ConnectSSH(addr, config.AppConfig.SSH.Username, config.AppConfig.SSH.Password, command)
	return rawOutput
}

func parseIntEthConfig(rawOutput string) (IntEthConfig map[string]models.EthInterface) {
	IntEthConfig = make(map[string]models.EthInterface)
	pattern := regexp.MustCompile(`(?m)^(interface.*)\n((?:.*\n)+?)!`)
	parsedOutput := pattern.FindAllStringSubmatch(rawOutput, -1)
	for _, matches := range parsedOutput {
		if len(matches) < 3 {
			log.Println("Could not parse EthInterfaces!")
		}
		interfaceName := matches[1]
		interfaceName = strings.TrimSpace(interfaceName)
		interfaceConfigString := matches[2]
		interfaceConfig := strings.Split(interfaceConfigString, "\r\n ")
		for i, configCommand := range interfaceConfig {
			interfaceConfig[i] = strings.TrimSpace(configCommand)
		}

		if strings.HasPrefix(interfaceName, "interface Vlan") {
			continue
		}
		IntEthConfig[interfaceName] = models.EthInterface{
			InterfaceConfig: interfaceConfig,
		}
	}
	return IntEthConfig
}

//TODO check if min 1 switch was reachable and sent data

func saveAsJson(switchInventory *[]models.NetworkSwitch) {
	_, err := ioUtil.ExistsDir("./database", true)
	if err != nil {
		log.Fatalln(err)
	}
	result, err := json.MarshalIndent(switchInventory, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}
	fileName := time.Now().Format("2006-01-02T15-04-05") + ".json"
	filePath := filepath.Join("./database", fileName)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = file.Write(result)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Successfuly saved database @%s", absFilePath)
}
