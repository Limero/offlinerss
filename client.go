package main

import (
	"errors"
	"os"

	"github.com/limero/offlinerss/models"
)

func (clients Clients) GetSyncToActions() ([]models.SyncToAction, error) {
	if len(clients) == 0 {
		return nil, errors.New("You have to enable at least one client in the config file")
	}

	// Grab changes from each client
	var syncToActions []models.SyncToAction

	for _, client := range clients {
		actions, err := client.GetChanges()
		if err != nil {
			return nil, err
		}
		syncToActions = append(syncToActions, actions...)
	}

	return syncToActions, nil
}

func (clients Clients) Sync(folders []*models.Folder) error {
	for _, client := range clients {
		if _, err := os.Stat(client.ReferenceDB()); errors.Is(err, os.ErrNotExist) {
			if err := client.CreateNewCache(); err != nil {
				return err
			}
		}

		if err := client.AddToCache(folders); err != nil {
			return err
		}
	}

	return nil
}
