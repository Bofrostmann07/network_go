package ethIntSearch

import (
	"fmt"
	"log"
	"net"
	"network_go/internal/SSHUtil"
	"network_go/internal/config"
	"regexp"
	"strings"
)

type NetworkSwitch struct {
	Address   string
	Hostname  string
	Platform  string
	Group     string
	Reachable bool

	IntEthConfig map[string][]string
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
	fmt.Println(switchInventory)
}

func sendSingleShowCommand(command string, networkSwitch NetworkSwitch) (rawOutput string) {
	addr := net.JoinHostPort(networkSwitch.Address, config.AppConfig.SSH.Port)
	rawOutput = SSHUtil.ConnectSSH(addr, config.AppConfig.SSH.Username, config.AppConfig.SSH.Password, command)
	return rawOutput
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
