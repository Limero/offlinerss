package main

import (
	"fmt"

	"github.com/limero/offlinerss/client/feedreader"
	"github.com/limero/offlinerss/client/newsboat"
	"github.com/limero/offlinerss/client/quiterss"
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
	"github.com/limero/offlinerss/server/miniflux"
	"github.com/limero/offlinerss/server/newsblur"
)

func run() error {
	config, err := getConfig()
	if err != nil {
		return err
	}

	var clients models.Clients
	for _, clientConfig := range config.Clients {
		switch clientConfig.Type {
		case "feedreader":
			clients = append(clients, feedreader.New(clientConfig))
		case "newsboat":
			clients = append(clients, newsboat.New(clientConfig))
		case "quiterss":
			clients = append(clients, quiterss.New(clientConfig))
		}
	}

	syncToActions, err := GetSyncToActions(clients)
	if err != nil {
		return err
	}

	var s models.Server
	switch config.Server.Type {
	case "miniflux":
		s = miniflux.New(config.Server)
	case "newsblur":
		s = newsblur.New(config.Server)
	}

	folders, err := SyncServer(s, syncToActions)
	if err != nil {
		return err
	}

	if err := SyncClients(clients, folders); err != nil {
		return err
	}

	log.Info("Everything synced!")

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println("Error:", err)
	}
}
