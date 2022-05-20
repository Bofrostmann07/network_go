package inventory

import (
	"encoding/csv"
	"fmt"
	"log"
	"network_go/internal/eth_int_search"
	"os"
)

func ReadSwitchInventoryFromCSV() []eth_int_search.NetworkSwitch {

	csvFile, err := os.Open("api/switchInventory.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	csvr := csv.NewReader(csvFile)
	csvr.FieldsPerRecord = -1
	fileData, err := csvr.ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	validatedFileData := validateCSVData(fileData)

	switchInventory := make([]eth_int_search.NetworkSwitch, len(fileData[1:]))
	for i, line := range validatedFileData[1:] {
		switchInvData := eth_int_search.NetworkSwitch{
			Hostname: line[0],
			Address:  line[1],
			Platform: line[2],
			Group:    line[3],
		}
		switchInventory[i] = switchInvData
	}
	return switchInventory
}

func validateCSVData(fileData [][]string) [][]string {
	errorCounter := 0
	for i, line := range fileData[1:] {
		if len(line) < 4 {
			errorCounter++
			log.Printf("%v/%v", "Too less arguments on line", i)
		}
		platform := line[2]
		if platform != "cisco_ios" && platform != "cisco_xe" {
			errorCounter++
			log.Printf("%v/%v", "Unsupported platform on line", i)
		}
	}
	if errorCounter > 0 {
		log.Fatalln("Invalid CSV data. Fix the listed issues to proceed with the program.")
	}
	log.Printf("%v/%v", "CSV data is valid. Found switches.", len(fileData)) //TODO Printf richten
	return fileData
}
