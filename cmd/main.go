package main

import "network_go/internal/eth_int_search"

func main() {
	appConfig := eth_int_search.AppConfig{
		Username: "automate",
		Password: "automateme",
		SSHPort:  "22",
	}
	networkSwitchList := []eth_int_search.NetworkSwitch{{Address: "10.20.0.1", Hostname: "test1"}, {Address: "10.20.0.1", Hostname: "test2"}}
	eth_int_search.FetchEthIntConfig(&networkSwitchList, appConfig)
	//maclookup.DoLookUp()
}
