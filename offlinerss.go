package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/limero/offlinerss/models"
	"github.com/mitchellh/go-homedir"
)

func getConfig() (*models.Config, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		configDir = filepath.Join(homeDir, ".config")
	}

	data, err := ioutil.ReadFile(filepath.Join(configDir, "offlinerss/config.json"))
	if err != nil {
		return nil, err
	}

	var config models.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Expand paths for each client (~ -> /home/username)
	// only i is used because it can overwrite the value
	for i := range config.Clients {
		if config.Clients[i].Paths.Cache, err = homedir.Expand(config.Clients[i].Paths.Cache); err != nil {
			return nil, err
		}

		if config.Clients[i].Paths.Urls, err = homedir.Expand(config.Clients[i].Paths.Urls); err != nil {
			return nil, err
		}
	}

	return &config, nil
}

func run() error {
	config, err := getConfig()
	if err != nil {
		return err
	}

	syncToActions, err := GetSyncToActions(config.Clients)
	if err != nil {
		return err
	}

	folders, err := DoSync(config.Server, syncToActions)
	if err != nil {
		return err
	}

	if err := GenerateDatabases(config.Clients, folders); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println("Error:", err)
	}
}
