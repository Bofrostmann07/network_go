package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
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
}

func ReadConfig() (Config, error) {
	conf, parseErr := parseConfigFromFile()
	if parseErr != nil {
		log.Fatalln("Error while parsing yaml file:", parseErr)
	}
	AppConfig = conf
	return conf, nil
}

func parseConfigFromFile() (Config, error) {
	yamlFileContent, readErr := ioutil.ReadFile(pathToConfigFile)
	if readErr != nil {
		return Config{}, readErr
	}

	conf, parseErr := parseConfigFromBytes(yamlFileContent)
	if parseErr != nil {
		return Config{}, parseErr
	}
	return conf, nil
}

func parseConfigFromBytes(data []byte) (Config, error) {
	var config Config

	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
