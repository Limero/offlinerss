package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/limero/offlinerss/client"
	"github.com/limero/offlinerss/models"
	"github.com/limero/offlinerss/server"
	_ "github.com/mattn/go-sqlite3"
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

	data, err := os.ReadFile(filepath.Join(configDir, "offlinerss/config.json"))
	if err != nil {
		return nil, err
	}

	var config models.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func run() error {
	config, err := getConfig()
	if err != nil {
		return err
	}

	var clients models.Clients
	for _, clientConfig := range config.Clients {
		switch clientConfig.Type {
		case "feedreader":
			clients = append(clients, client.NewFeedreader(clientConfig))
		case "newsboat":
			clients = append(clients, client.NewNewsboat(clientConfig))
		case "quiterss":
			clients = append(clients, client.NewQuiteRSS(clientConfig))
		}
	}

	syncToActions, err := GetSyncToActions(clients)
	if err != nil {
		return err
	}

	var s models.Server
	switch config.Server.Type {
	case "miniflux":
		s = server.NewMiniflux(config.Server)
	case "newsblur":
		s = server.NewNewsblur(config.Server)
	}

	folders, err := SyncServer(s, syncToActions)
	if err != nil {
		return err
	}

	if err := SyncClients(clients, folders); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println("Error:", err)
	}
}
