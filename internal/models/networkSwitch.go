package models

type NetworkSwitch struct {
	Address   string `json:"Address"`
	Hostname  string `json:"Hostname"`
	Platform  string `json:"Platform"`
	Group     string `json:"Group"`
	Reachable bool   `json:"Reachable"`

	EthInterfaces map[string]EthInterface `json:"EthInterfaces"`
}

type EthInterface struct {
	InterfaceConfig []string `json:"InterfaceConfig"`
	MacList         []string `json:"MacList"`
}
