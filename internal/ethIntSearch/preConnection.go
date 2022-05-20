package ethIntSearch

import (
	"fmt"
	"log"
	"net"
	"network_go/internal/config"
	"time"
)

func CheckConnection(switchInventory *[]NetworkSwitch) {
	checkConnectivity(switchInventory)
}

func checkConnectivity(switchInventory *[]NetworkSwitch) {
	errorCounter := 0
	for i, switchElement := range *switchInventory {
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(switchElement.Address, config.AppConfig.SSH.Port), timeout)
		if err != nil {
			errorCounter++
			log.Println("Connecting error:", err)
		}
		if conn != nil {
			defer conn.Close()
			(*switchInventory)[i].Reachable = true
		}
	}
	if errorCounter != 0 {
		log.Printf("Couldn't reach %d switches. All unreachable switches will be left out.", errorCounter)

		var userInput string
		fmt.Print("Continue? [y]/n: ")
		//_, err := fmt.Scanf("\n%s", &userInput) //TODO unexpected newline
		//if err != nil {
		//	fmt.Print(err)
		//}
		fmt.Scanln(&userInput)
		if userInput == "\n" || userInput == "y" {
			log.Println("Continuing with the connection check")
			return
		}
		log.Fatalln("Stopping. Restart the tool after fixing the connection issues.")
	}
}
