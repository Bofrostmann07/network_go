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
	"network_go/internal/parser"
	"network_go/internal/util/ioUtil"
	"os"
	"path/filepath"
	"sort"
	"time"
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

	fmt.Print("Query: ")
	searchQuery := ioUtil.ReadLine()

	matchedSwitches, notMatchedSwitches := querySearch(searchQuery, &switchInventory)

	fmt.Print("Filter: ")
	filterQuery := ioUtil.ReadLine()

	filteredNetworkSwitches, notFilterMatchedSwitches := filterInterfaces(filterQuery, &matchedSwitches)
	notMatchedSwitches = append(notMatchedSwitches, notFilterMatchedSwitches...)

	fmt.Printf("Matched %d switches.", len(filteredNetworkSwitches))
	fmt.Printf("No interfaces matched for %d switches.\n", len(notMatchedSwitches))

	asaveAsJson(filteredNetworkSwitches, notMatchedSwitches)
}

func getDatabaseData() []models.NetworkSwitch {
	sortedFileList, err := getFileListDesc()
	if err != nil {
		switchInventory := inventory.ReadSwitchInventoryFromCSV()
		fetchEthIntConfig(&switchInventory)
		return switchInventory
	}
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
		return switchInventory
	}

	log.Println("List of 10 recent files @./database:")
	for i := 0; i < 10; i++ {
		log.Println(sortedFileList[i].Name())
	}
	fmt.Print("Select file: ")
	selectedFile := ioUtil.ReadLine()
	switchInventory := readDatabase(selectedFile)
	return switchInventory
}

func getFileListDesc() ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir("./database")
	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		return []fs.FileInfo{}, errors.New("no files found")
	}

	// descending order
	sort.Slice(files,
		func(i, j int) bool {
			return files[i].ModTime().After(files[j].ModTime())
		})
	return files, nil
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

// querySearch searches the switch inventory for the given search query.
func querySearch(query string, switchinventory *[]models.NetworkSwitch) (matched, notMatched []models.NetworkSwitch) {
	parsedQuery, err := parser.ParseQuery(query)
	if err != nil {
		log.Fatalln("Query not parsable. ", err)
	}

	for _, networkSwitch := range *switchinventory {
		matches := networkSwitch.EvaluateQuery(parsedQuery)
		if matches {
			matched = append(matched, networkSwitch)
		} else {
			notMatched = append(notMatched, networkSwitch)
		}
	}
	fmt.Println(matched)
	fmt.Printf("Found %d switches.\n", len(matched))

	return matched, notMatched
}

func filterInterfaces(query string, switchinventory *[]models.NetworkSwitch) (matchedSwitches, noMatch []models.NetworkSwitch) {
	parsedQuery, err := parser.ParseQuery(query)
	if err != nil {
		log.Fatalln("Query not parsable. ", err)
	}

	for _, networkSwitch := range *switchinventory {
		newInterfaces := make(map[string]models.EthInterface)
		// For each network interface, check if the interface matches the query.
		for key, networkInterface := range networkSwitch.EthInterfaces {
			matches := networkInterface.EvaluateQuery(parsedQuery)
			if matches {
				newInterfaces[key] = networkInterface
			}
		}
		networkSwitch.EthInterfaces = newInterfaces

		if len(newInterfaces) > 0 {
			matchedSwitches = append(matchedSwitches, networkSwitch)
		} else {
			noMatch = append(noMatch, networkSwitch)
		}
	}

	fmt.Println(matchedSwitches)
	return matchedSwitches, noMatch
}

func asaveAsJson(filteredSwitchList []models.NetworkSwitch, notMatchedSwitchList []models.NetworkSwitch) {
	_, err := ioUtil.ExistsDir("./results", true)
	if err != nil {
		log.Fatalln(err)
	}
	result, err := json.MarshalIndent(filteredSwitchList, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}

	result2, err := json.MarshalIndent(notMatchedSwitchList, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}

	fileName := time.Now().Format("2006-01-02T15-04-05") + ".json"
	filePath := filepath.Join("./results", fileName)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = file.Write([]byte("Matched switches:\n"))
	if err != nil {
		log.Fatalln(err)
	}

	_, err = file.Write(result)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = file.Write([]byte("\n\n\nNot matched switches:\n"))
	if err != nil {
		log.Fatalln(err)
	}

	_, err = file.Write(result2)
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Successfuly saved database @%s", absFilePath)
}
