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
	Target               string //The hostname or IP of the target for enumeration.
	FileRoot             string //The root configuration file directory
	GlobalConfigFileName string
	OutputDirectory      string //The base path to the directory where treenum's output is stored.
	ScriptDirName        string //The user-defined scan script dir, which houses their service scan scripts
}

func (config *Config) GetScriptDir() string {
	return config.GetFileRoot() + "/scripts/" + config.GetScriptDirName()
}

func (config *Config) GetScriptDirName() string {
	if config.ScriptDirName == "" {
		config.ScriptDirName = "default"
	}
	return config.ScriptDirName
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

func (config *Config) GetFileRoot() string {
	if config.FileRoot == "" {
		config.FileRoot = os.Getenv("HOME") + "/.config/treenum"
	}
	return config.FileRoot
}

func (config *Config) GetGlobalConfigFilename() string {
	if config.GlobalConfigFileName == "" {
		config.GlobalConfigFileName = "config.json"
	}
	return config.GlobalConfigFileName
}

func (config *Config) GetConfigFilePath() string {
	return config.GetFileRoot() + "/" + config.GetGlobalConfigFilename()
}

/* Parses the configuration file into a Config object, then returns it.
 * Can also either generate an error in attempting to read the file, or unmarshaling it
 * into a json object. */
func (config *Config) Load() error {
	var err error
	b, err := ioutil.ReadFile(config.GetConfigFilePath())
	if err == nil {
		err = json.Unmarshal(b, &config)
	}
	return err
}
