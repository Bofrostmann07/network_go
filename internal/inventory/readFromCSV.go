package inventory

import (
	"encoding/csv"
	"fmt"
	"log"
	"network_go/internal/ethIntSearch"
	"os"
)

func ReadSwitchInventoryFromCSV() []ethIntSearch.NetworkSwitch {

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

	switchInventory := make([]ethIntSearch.NetworkSwitch, len(fileData[1:]))
	for i, line := range validatedFileData[1:] {
		switchInvData := ethIntSearch.NetworkSwitch{
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
	log.Printf("CSV data is valid. Found %d switches.", len(fileData))
	return fileData
}

//TODO check if csv file has at least one entry
