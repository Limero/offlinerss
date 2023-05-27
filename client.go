package main

import (
	"errors"
	"os"

	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
)

func GetSyncToActions(clients models.Clients) (models.SyncToActions, error) {
	if len(clients) == 0 {
		return nil, errors.New("You have to enable at least one client in the config file")
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
			log.Info("  ☁️  %d unstarred", unstarred)
		}

		syncToActions = append(syncToActions, actions...)
	}

	if len(syncToActions) == 0 {
		log.Info("No changes found for any local clients")
	}

	return syncToActions, nil
}

func SyncClients(clients models.Clients, folders []*models.Folder) error {
	for _, client := range clients {
		if _, err := os.Stat(client.ReferenceDB()); errors.Is(err, os.ErrNotExist) {
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
