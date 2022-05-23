package ethIntSearch

import (
	"log"
	"net"
	"network_go/internal/config"
	"network_go/internal/models"
	"network_go/internal/util/ioUtil"
	"time"
)

func CheckConnection(switchInventory *[]models.NetworkSwitch) {
	jobs := make(chan models.NetworkSwitch, len(*switchInventory))
	results := make(chan models.NetworkSwitch, len(*switchInventory))

	for i := 0; i < config.AppConfig.SSH.MaxGoRoutines; i++ {
		go worker(jobs, results)
	}
	for _, networkSwitch := range *switchInventory {
		jobs <- networkSwitch
	}
	close(jobs)

	for i := 0; i < len(*switchInventory); i++ {
		(*switchInventory)[i] = <-results
		//fmt.Println(<-results)
	}
	close(results)

	userInputIfSwitchesUnreachable(switchInventory)
}

func worker(jobs <-chan models.NetworkSwitch, results chan<- models.NetworkSwitch) {
	for networkSwitch := range jobs {
		results <- checkConnectivity(networkSwitch)
	}
}

func checkConnectivity(networkSwitch models.NetworkSwitch) models.NetworkSwitch {
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
}

func userInputIfSwitchesUnreachable(switchInventory *[]models.NetworkSwitch) {
	unreachableCounter := 0
	for _, networkSwitch := range *switchInventory {
		if !networkSwitch.Reachable {
			unreachableCounter++
		}
	}

	if unreachableCounter != 0 {
		log.Printf("Couldn't reach %d switches. All unreachable switches will be left out.", unreachableCounter)
		userAnswer := ioUtil.UserInputYesNo("Continue? [y]/n: ", true)
		if userAnswer {
			log.Println("Continuing with the connection check")
			return
		}
		log.Fatalln("Stopping. Restart the tool after fixing the connection issues.")
	}
}
