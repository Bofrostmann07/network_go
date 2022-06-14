package inventory

import (
	"encoding/csv"
	"fmt"
	"log"
	"network_go/internal/models"
	"os"
)

func ReadSwitchInventoryFromCSV() []models.NetworkSwitch {

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

	switchInventory := make([]models.NetworkSwitch, len(fileData[1:]))
	for i, line := range validatedFileData[1:] {
		switchInvData := models.NetworkSwitch{
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
			log.Printf("Too less arguments on line %d", i)
		}
		platform := line[2]
		if platform != "cisco_ios" && platform != "cisco_xe" {
			errorCounter++
			log.Printf("Unsupported platform on line %d", i)
		}
	}
	if errorCounter > 0 {
		log.Fatalln("Invalid CSV data. Fix the listed issues to proceed with the program.")
	}
	log.Printf("CSV data is valid. Found %d switches.", len(fileData[1:]))
	return fileData
}

//TODO check if csv file has at least one entry
