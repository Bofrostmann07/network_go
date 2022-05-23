package ethIntSearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"network_go/internal/inventory"
	"network_go/internal/models"
	"network_go/internal/util/ioUtil"
	"os"
	"sort"
)

func SwitchSearch() {
	log.Println("Started switch search")
	switchInventory := getDatabaseData()
	log.Printf("Loaded %d switches", len(switchInventory))

	fmt.Printf("Usage: \"[search mask]\" [flags]\n" +
		"Flags/Options:\n" +
		"--n     Turn to negative search mode. Will list all interfaces, which wont fit the search mask.\n" +
		"--o     Tries to strip off out-of-band-management interfaces\n" +
		"--u     Tries to strip off uplink interfaces\n" +
		"Example: \"switchport mode access\" --n --u\n")

}

func getDatabaseData() []models.NetworkSwitch {
	sortedFileList := getFileListDesc()
	recentFile := "./database/" + sortedFileList[0].Name()
	log.Printf("The most recent Datafile is from %s.", recentFile)
	loadRecentFile := ioUtil.UserInputYesNo("Open most recent file? [y]/n: ", true)
	if loadRecentFile {
		switchInventory := readDatabase(recentFile)
		return switchInventory
	}

	retrieveNow := ioUtil.UserInputYesNo("Retrieve switch config now? [y]/n: ", true)
	if retrieveNow {
		switchInventory := inventory.ReadSwitchInventoryFromCSV()
		fetchEthIntConfig(&switchInventory)
		return nil //TODO this is wrong...
	}

	log.Println("List of 10 recent files @./database:")
	for i := 0; i < 10; i++ {
		log.Println(sortedFileList[i].Name())
	}
	var selectedFile string
	fmt.Print("Select file: ")
	_, err := fmt.Scanln(&selectedFile)
	if err != nil {
		fmt.Print(err)
	}
	switchInventory := readDatabase(selectedFile)
	return switchInventory
}

func getFileListDesc() []fs.FileInfo {
	files, err := ioutil.ReadDir("./database")
	if err != nil {
		log.Fatal(err)
	}

	// descending order
	sort.Slice(files,
		func(i, j int) bool {
			return files[i].ModTime().After(files[j].ModTime())
		})
	return files
}

func readDatabase(file string) []models.NetworkSwitch {
	fileData, err := ioutil.ReadFile(file)
	if errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Could not open Database %q", file)
	}
	var switchInventory []models.NetworkSwitch
	err = json.Unmarshal(fileData, &switchInventory)
	if err != nil {
		log.Fatalln(err)
	}
	return switchInventory
}

func saveSearchResult() {

}
