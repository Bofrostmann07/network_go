package macLookup

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type searchedMac struct {
	sourceMac    string
	formattedMac string
	formattedOui string
	vendorFound  bool
	vendor       string
}

var pathIeeeRegister string = "internal/macLookup/IEEE_MAC_register.json"

func DoLookUp() {
	log.Println("Started MAC address lookup tool")
	ieeeRegistry := checkIeeeRegistry()
	InputText := userInput()
	parsedMACAddressList := parseMacAddresses(InputText)
	searchQuery := formatMACAddresses(parsedMACAddressList)
	searchIeeeRegistry(&searchQuery, ieeeRegistry)
	printResult(&searchQuery)
}

//Download IEEE MAC registry if it`s missing or outdated
func checkIeeeRegistry() (ieeeRegistry map[string]string) {
	file, err := os.OpenFile(pathIeeeRegister, os.O_RDWR, 0755)
	if os.IsNotExist(err) {
		log.Println("IEEE MAC Registry is missing. Downloading...")
		getIeeeMacList()
		return ieeeRegistry
	}
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	fileInfo, err := os.Stat(pathIeeeRegister)
	if err != nil {
		fmt.Println(err)
	}
	modifiedTime := fileInfo.ModTime()
	expiryTime := modifiedTime.AddDate(0, 0, 7)
	currentTime := time.Now()
	//fmt.Println("Last modified time : ", modifiedTime)
	//fmt.Println("Expiry time : ", expiryTime)
	//fmt.Println("Current time : ", currentTime)

	if currentTime.After(expiryTime) {
		log.Println("IEEE MAC Registry is expired. Downloading...")
		getIeeeMacList()
		return ieeeRegistry
	}

	log.Println("IEEE MAC Registry is valid. [1/]")
	ieeeRegistry = readIeeeRegistry()
	return ieeeRegistry
}

//Multi-line user input
func userInput() (lines []string) {
	fmt.Println("Complete input by entering 'exit' in new line.\nEnter Lines:")
	scn := bufio.NewScanner(os.Stdin)
	for {
		for scn.Scan() {
			line := scn.Text()
			// Group Separator (GS ^]): ctrl-]
			if len(line) == 1 && line[0] == '\x1D' {
				break
			}
			if line == "exit" {
				break
			}
			lines = append(lines, line)
		}

		if err := scn.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}
		if len(lines) == 0 {
			fmt.Println("Input is empty")
			os.Exit(0)
		}
		break
	}
	return lines
}

//Parse out all the MAC addresses from user input
func parseMacAddresses(InputText []string) (parsedMACAddressList []string) {
	pattern := regexp.MustCompile(`(?:[[:xdigit:]]{2}[-:.]){5}[[:xdigit:]]{2}|(?:[[:xdigit:]]{4}.){2}[[:xdigit:]]{4}|[[:xdigit:]]{12}`)

	for _, line := range InputText {
		MacAddress := pattern.FindAllString(line, -1)
		parsedMACAddressList = append(parsedMACAddressList, MacAddress...)
	}

	parsedMACAddressList = removeDuplicateValues(parsedMACAddressList)

	return parsedMACAddressList
}

//Remove all duplicate MAC addresses out of user input list
func removeDuplicateValues(inputSlice []string) (uniqueSlice []string) {
	keys := make(map[string]bool)

	for _, entry := range inputSlice {
		if _, ok := keys[entry]; ok {
			continue
		}
		keys[entry] = true
		uniqueSlice = append(uniqueSlice, entry)
	}
	log.Println("Parsed", len(uniqueSlice), "unique MAC addresses.")
	return uniqueSlice
}

//Download IEEE MAC list and format it to JSON
func getIeeeMacList() (ieeeRegistry map[string]string) {
	//Download IEEE MAC list
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://standards-oui.ieee.org/oui/oui.txt", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%s\n", bodyText)

	//Format IEEE MAC list to JSON oui:vendor
	ieeeRegistry = make(map[string]string)
	pattern := regexp.MustCompile(`([[:xdigit:]]{6})\s+\(.*\)\s+(.*)`)
	output := pattern.FindAllStringSubmatch(string(bodyText), -1)
	//fmt.Println(output)
	for _, element := range output {
		formattedVendor := strings.TrimSuffix(element[2], "\r")
		ieeeRegistry[element[1]] = formattedVendor
	}

	jsonString, _ := json.Marshal(ieeeRegistry)
	if err != nil {
		log.Fatal(err)
	}

	//Write IEEE MAC list as JSON
	file, err := os.Create(pathIeeeRegister)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err2 := file.Write(jsonString)
	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("Updated IEEE MAC address registry. [1/]")
	return ieeeRegistry
}

//Open and read IEEE MAC list
func readIeeeRegistry() (ieeeRegistry map[string]string) {
	jsonFile, err := os.Open(pathIeeeRegister)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &ieeeRegistry)
	return ieeeRegistry
}

func formatMACAddresses(parsedMACAddressList []string) (searchQuery []searchedMac) {
	searchQuery = make([]searchedMac, len(parsedMACAddressList))
	pattern := regexp.MustCompile(`[-:.]`)
	for i, element := range parsedMACAddressList {
		searchQuery[i].sourceMac = element
		element = pattern.ReplaceAllString(element, "")
		element = strings.ToUpper(element)
		searchQuery[i].formattedMac = element
		searchQuery[i].formattedOui = element[0:6]
	}
	return
}

func searchIeeeRegistry(searchQuery *[]searchedMac, ieeeRegistry map[string]string) {
	foundEntries := 0
	for oui, registryEntry := range ieeeRegistry {
		for i, searchElement := range *searchQuery {
			if searchElement.formattedOui == oui {
				(*searchQuery)[i].vendor = registryEntry
				(*searchQuery)[i].vendorFound = true
				foundEntries++
			}
		}
	}
	log.Printf("Found %d addresses in IEEE MAC Registry.\n", foundEntries)
}

func printResult(searchQuery *[]searchedMac) {
	for i, element := range *searchQuery {
		log.Printf("%d - %s: %s", i, element.sourceMac, element.vendor)
	}
}
