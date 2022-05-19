package eth_int_search

import (
	"fmt"
	"log"
	"network_go/internal/sshUtil"
	"regexp"
	"strings"
)

type AppConfig struct {
	Username string
	Password string
	SSHPort  string
}

var appConfig AppConfig

type NetworkSwitch struct {
	Address   string
	Hostname  string
	Platform  string
	Reachable bool

	IntEthConfig map[string][]string
}

func sendSingleShowCommand(command string, networkSwitch NetworkSwitch) (rawOutput string) {
	addr := fmt.Sprintf("%s:%s", networkSwitch.Address, appConfig.SSHPort)
	rawOutput = sshUtil.ConnectSSH(addr, appConfig.Username, appConfig.Password, command)
	return rawOutput
}

func FetchEthIntConfig(switchInventory *[]NetworkSwitch, config AppConfig) {
	appConfig = config
	if switchInventory == nil {
		return
	}

	for i, networkSwitch := range *switchInventory {
		rawOutput := sendSingleShowCommand("show derived-config | begin interface", networkSwitch)
		(*switchInventory)[i].IntEthConfig = parseIntEthConfig(rawOutput)
	}
	fmt.Println(switchInventory)
}

func parseIntEthConfig(rawOutput string) (IntEthConfig map[string][]string) {
	IntEthConfig = make(map[string][]string)
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
		IntEthConfig[interfaceName] = interfaceConfig
	}
	return IntEthConfig
}
