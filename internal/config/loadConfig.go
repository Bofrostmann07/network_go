package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"network_go/internal/util/ioUtil"
	"os"
)

var pathToConfigFile = "configs/appConfig.yml"
var AppConfig Config

type Config struct {
	InputSource string `yaml:"input_source"`
	DebugMode   bool   `yaml:"debug_mode"`
	SSH         struct {
		Username                string `yaml:"ssh_username"`
		Password                string `yaml:"ssh_password"`
		Port                    string `yaml:"ssh_port"`
		Timeout                 int    `yaml:"ssh_timeout"`
		MaxGoRoutines           int    `yaml:"max_threads"`
		SkipReachabilityCheck   bool   `yaml:"skip_ssh_reachability_check"`
		SkipAuthenticationCheck bool   `yaml:"skip_ssh_authentication_check"`
	} `yaml:"ssh_config"`
	Prime struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Address  string `yaml:"address"`
		Group    string `yaml:"group"`
	} `yaml:"prime_config"`
	MacAddressLookup struct {
		ApiToken string `yaml:"api_token_macvendors"`
	} `yaml:"mac_address_lookup"`
	SystemValues struct {
		AppConfigVersion string `yaml:"app_config_version"`
	} `yaml:"system_values"`
}

func ReadConfig() {
	// Create directory 'configs', if it's not existing
	_, err := ioUtil.ExistsDir("./configs", true)
	if err != nil {
		log.Fatalln("Could not create configs directory.")
	}

	getConfigFromFile()
}

func getConfigFromFile() {
	appConfigFile, readErr := os.ReadFile(pathToConfigFile)
	if os.IsNotExist(readErr) {
		log.Fatalln("Error while reading yaml file: ", readErr)
	}
	if readErr != nil {
		log.Fatalln("Error while reading yaml file: ", readErr)
	}

	parseErr := yaml.Unmarshal(appConfigFile, &AppConfig)
	if parseErr != nil {
		log.Fatalln("Error while unmarshalling yaml file: ", parseErr)
	}
}



