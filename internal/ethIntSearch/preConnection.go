package ethIntSearch

import (
	"fmt"
	"log"
	"net"
	"network_go/internal/config"
	"time"
)

func CheckConnection(switchInventory *[]NetworkSwitch) {
	jobs := make(chan NetworkSwitch, len(*switchInventory))
	results := make(chan NetworkSwitch, len(*switchInventory))

	for i := 0; i < config.AppConfig.SSH.MaxGoRoutines; i++ {
		go worker(jobs, results)
	}
	for _, networkSwitch := range *switchInventory {
		jobs <- networkSwitch
	}
	//close(jobs)

	for i := 0; i < len(*switchInventory); i++ {
		(*switchInventory)[i] = <-results
		fmt.Println(<-results)

	}
}

func worker(jobs <-chan NetworkSwitch, results chan<- NetworkSwitch) {
	for networkSwitch := range jobs {
		results <- checkConnectivity(networkSwitch)
	}
}

func checkConnectivity(networkSwitch NetworkSwitch) NetworkSwitch {
	timeout := time.Duration(config.AppConfig.SSH.Timeout) * time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(networkSwitch.Address, config.AppConfig.SSH.Port), timeout)
	if err != nil {
		log.Println("Connecting error:", err)
	}
	if conn != nil {
		defer conn.Close()
		networkSwitch.Reachable = true
	}
	return networkSwitch
	//if errorCounter != 0 {
	//	log.Printf("Couldn't reach %d switches. All unreachable switches will be left out.", errorCounter)
	//
	//	var userInput string
	//	fmt.Print("Continue? [y]/n: ")
	//	//_, err := fmt.Scanf("\n%s", &userInput) //TODO unexpected newline
	//	//if err != nil {
	//	//	fmt.Print(err)
	//	//}
	//	fmt.Scanln(&userInput)
	//	if userInput == "\n" || userInput == "y" {
	//		log.Println("Continuing with the connection check")
	//		return
	//	}
	//	log.Fatalln("Stopping. Restart the tool after fixing the connection issues.")
}
