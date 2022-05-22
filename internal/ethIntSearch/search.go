package ethIntSearch

import (
	"encoding/json"
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"network_go/internal/util/ioUtil"
	"os"
	"sort"
)

func SwitchSearch() {
	log.Println("Started switch search")
	getDatabaseData()
}

func getDatabaseData() {
	sortedFileList := getFileListDesc()
	recentFile := "./database/" + sortedFileList[0].Name()
	log.Printf("The most recent Datafile is from %s.", recentFile)
	loadRecentFile := ioUtil.UserInputYesNo("Open most recent file: [y]/n", true)
	if loadRecentFile {
		readDatabase(recentFile)
	}
	//retrieveNow := ioUtil.UserInputYesNo("Retrieve switch config now: [y]/n", true)
	//if retrieveNow {
	//	switchInventory := inventory.ReadSwitchInventoryFromCSV()
	//	FetchEthIntConfig(&switchInventory)
	//}
	//log.Println("List of 10 recent files @./database:")
	//for i := 0; i < 10; i++ {
	//	log.Println(sortedFileList[i].Name())
	//}
	//var selectedFile string
	//fmt.Print("Select file: ")
	//_, err := fmt.Scanln(&selectedFile)
	//if err != nil {
	//	fmt.Print(err)
	//}
	//readDatabase(selectedFile)
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

func readDatabase(file string) {
	fileData, err := ioutil.ReadFile(file)
	if errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Could not open Database %q", file)
	}
	err = json.Unmarshal(fileData, &NetworkSwitch{})
	if err != nil {
		log.Fatalln(err)
	}
}

func saveSearchResult() {

}
