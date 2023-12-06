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
		switch clientConfig.Name {
		case models.ClientFeedReader:
			clients = append(clients, feedreader.New(clientConfig))
		case models.ClientNewsboat:
			clients = append(clients, newsboat.New(clientConfig))
		case models.ClientQuiteRSS:
			clients = append(clients, quiterss.New(clientConfig))
		}
	}

	syncToActions, err := GetSyncToActions(clients)
	if err != nil {
		return err
	}

	var s models.Server
	switch config.Server.Name {
	case models.ServerMiniflux:
		s = miniflux.New(config.Server)
	case models.ServerNewsBlur:
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
