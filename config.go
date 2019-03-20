package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

/* Holds all of the configuration settings used to modify AutoEnum's behavior.
   Should be read in using getConfig(), and modified based on command line arguments (cmd taking precedence) */
type Config struct {
	Target          string //The hostname or IP of the target for enumeration.
	OutputDirectory string //The base path to the directory where treenum's output is stored.
}

func (config *Config) GetOutputDirectory() string {
	if config.OutputDirectory == "" {
		var err error
		config.OutputDirectory, err = os.Getwd()
		if err != nil {
			log.Println("Warning: Could not retrieve working directory.")
			log.Println(err)
			log.Println("Setting output directory to /tmp/treenum")
			config.OutputDirectory = "/tmp/treenum"
		}
	}
	return config.OutputDirectory
}

var configFileDir = os.Getenv("HOME") + "/.config/treenum"

const configFilename = "config.json"

/* Parses the configuration file into a Config object, then returns it.
 * Can also either generate an error in attempting to read the file, or unmarshaling it
 * into a json object. */
func getConfig() (Config, error) {
	var ret Config
	b, err := ioutil.ReadFile(configFileDir + "/" + configFilename)
	if err == nil {
		err = json.Unmarshal(b, &ret)
	}
	return ret, err
}
