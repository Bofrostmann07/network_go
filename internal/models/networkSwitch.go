package models

import (
	"errors"
	"fmt"
	"network_go/internal/parser/ast"
	"strconv"
)

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

func (e EthInterface) Search(field ast.Field) bool {
	switch field.Bucket {
	case "interface":
		return field.SearchStringSlice(e.InterfaceConfig)
	case "mac":
		return field.SearchStringSlice(e.MacList)
	}
	return false
}

func (n NetworkSwitch) Search(field ast.Field) bool {
	var searchString string

	switch field.Bucket {
	case "address":
		searchString = n.Address
	case "hostname":
		searchString = n.Hostname
	case "platform":
		searchString = n.Platform
	case "group":
		searchString = n.Group
	case "reachable":
		searchString = strconv.FormatBool(n.Reachable)
	case "interface", "mac":
		newEthInterfaces := make(map[string]EthInterface)

		for s, ethInterface := range n.EthInterfaces {
			if ethInterface.Search(field) {
				newEthInterfaces[s] = ethInterface
			}
		}
		if len(newEthInterfaces) == 0 {
			return false
		}

		//n.EthInterfaces = newEthInterfaces
		return true
	}

	return field.SearchString(searchString)
}

func (n NetworkSwitch) filterInterfaceData(field ast.Field) (*NetworkSwitch, error) {
	newEthInterfaces := make(map[string]EthInterface)

	for s, ethInterface := range n.EthInterfaces {
		if field.SearchStringSlice(ethInterface.MacList) {
			newEthInterfaces[s] = ethInterface
		}
	}
	if len(newEthInterfaces) == 0 {
		return nil, errors.New("no match found")
	}

	n.EthInterfaces = newEthInterfaces
	return &n, nil
}

func (n NetworkSwitch) EvaluateQuery(query *ast.Query) (matches bool) {
	matches = n.Search(*query.Field)

	// TODO handle nil query

	// No additional queries - Field only
	if query.AndQuery == nil && query.OrQuery == nil {
		return matches
	}

	// Both Queries
	if query.AndQuery != nil && query.OrQuery != nil {
		return (matches && n.EvaluateQuery(query.AndQuery)) || n.EvaluateQuery(query.OrQuery)
	}

	// And Query
	if query.AndQuery != nil {
		return matches && n.EvaluateQuery(query.AndQuery)
	}

	// Or Query
	if query.OrQuery != nil {
		return matches || n.EvaluateQuery(query.OrQuery)
	}

	fmt.Printf("Uncaught state - %v\n", query)
	return false
}
