package main

import (
	"fmt"
	"network_go/internal/config"
	"network_go/internal/ethIntSearch"
	"network_go/internal/macLookup"
	"network_go/internal/util/ioUtil"
)

func main() {
	fmt.Println("Started network_go\n")
	menu()
}

func menu() {
	_, err := config.ReadConfig()
	if err != nil {
		return
	}

	for true {
		fmt.Println("Please choose the Tool by number:\n" +
			"1 - Interface search\n" +
			"2 - MAC address batch lookup\n" +
			"99 - Show Config Values (global_config.yml)",
		)

		//var toolNumber int
		fmt.Print("Tool number: ")
		//_, err := fmt.Scanln(&toolNumber)
		toolNumber := ioUtil.ReadLine()

		switch toolNumber {
		case "1":
			ethIntSearch.SwitchSearch()
		case "2":
			macLookup.DoLookUp()
		case "99":
			println("N/A")
		default:
			println("Invalid tool number.")
		}
	}
}
