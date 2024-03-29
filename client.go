package main

import (
	"errors"

	"github.com/limero/offlinerss/client/feedreader"
	"github.com/limero/offlinerss/client/newsboat"
	"github.com/limero/offlinerss/client/quiterss"
	"github.com/limero/offlinerss/helpers"
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
)

func getClients(clientConfigs []models.ClientConfig) models.Clients {
	clients := make(models.Clients, len(clientConfigs))
	for i, clientConfig := range clientConfigs {
		switch clientConfig.Name {
		case models.ClientFeedReader:
			clients[i] = feedreader.New(clientConfig)
		case models.ClientNewsboat:
			clients[i] = newsboat.New(clientConfig)
		case models.ClientQuiteRSS:
			clients[i] = quiterss.New(clientConfig)
		default:
			panic("unknown client " + clientConfig.Name)
		}
	}
	return clients
}

func GetSyncToActions(clients models.Clients) (models.SyncToActions, error) {
	if len(clients) == 0 {
		return nil, errors.New("you have to enable at least one client in the config file")
	}

	// Grab changes from each client
	var syncToActions models.SyncToActions

	for _, client := range clients {
		actions, err := client.GetChanges()
		if err != nil {
			return nil, err
		}

		if len(actions) == 0 {
			continue
		}

		log.Info("Found %d changes in %s", len(actions), client.Name())
		read, unread, starred, unstarred := actions.SumActionTypes()
		if read > 0 {
			log.Info("  📖 %d read", read)
		}
		if unread > 0 {
			log.Info("  📕 %d unread", unread)
		}
		if starred > 0 {
			log.Info("  ⭐ %d starred", starred)
		}
		if unstarred > 0 {
			log.Info("  ☁️ %d unstarred", unstarred)
		}

		syncToActions = append(syncToActions, actions...)
	}

	if len(syncToActions) == 0 {
		log.Info("No changes found for any local clients")
	}

	return syncToActions, nil
}

func SyncClients(clients models.Clients, folders models.Folders) error {
	for _, client := range clients {
		if !helpers.FileExists(client.ReferenceDB()) {
			if err := client.CreateNewCache(); err != nil {
				return err
			}
		}

		log.Info("Syncing stories to local client %s", client.Name())
		if err := client.AddToCache(folders); err != nil {
			return err
		}
	}

	return nil
}
