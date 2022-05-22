package ethIntSearch

import (
	"encoding/json"
	"log"
	"net"
	"network_go/internal/config"
	"network_go/internal/util/ioUtil"
	"network_go/internal/util/sshUtil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type NetworkSwitch struct {
	Address   string
	Hostname  string
	Platform  string
	Group     string
	Reachable bool

	IntEthConfig map[string]EthInterface
}

type EthInterface struct {
	InterfaceConfig []string
	macList         []string
}

func FetchEthIntConfig(switchInventory *[]NetworkSwitch) {
	CheckConnection(switchInventory)
	if switchInventory == nil {
		return
	}

	for i, networkSwitch := range *switchInventory {
		if networkSwitch.Reachable {
			rawOutput := sendSingleShowCommand("show derived-config | begin interface", networkSwitch)
			(*switchInventory)[i].IntEthConfig = parseIntEthConfig(rawOutput)
		}
	}
	saveAsJson(switchInventory)
}

func sendSingleShowCommand(command string, networkSwitch NetworkSwitch) (rawOutput string) {
	addr := net.JoinHostPort(networkSwitch.Address, config.AppConfig.SSH.Port)
	rawOutput = sshUtil.ConnectSSH(addr, config.AppConfig.SSH.Username, config.AppConfig.SSH.Password, command)
	return rawOutput
}

func parseIntEthConfig(rawOutput string) (IntEthConfig map[string]EthInterface) {
	IntEthConfig = make(map[string]EthInterface)
	pattern := regexp.MustCompile(`(?m)^(interface.*)\n((?:.*\n)+?)!`)
	parsedOutput := pattern.FindAllStringSubmatch(rawOutput, -1)
	for _, matches := range parsedOutput {
		if len(matches) < 3 {
			log.Println("Could not parse IntEthConfig!")
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
		IntEthConfig[interfaceName] = EthInterface{
			InterfaceConfig: interfaceConfig,
		}
	}
	return IntEthConfig
}

//TODO check if min 1 switch was reachable and sent data

func saveAsJson(switchInventory *[]NetworkSwitch) {
	_, err := ioUtil.ExistsDir("./database", true)
	if err != nil {
		log.Fatalln(err)
	}
	result, err := json.MarshalIndent(switchInventory, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}
	fileName := time.Now().Format("2006-01-02T15-04-05 Mon") + ".json"
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
}
