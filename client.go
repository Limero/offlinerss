package main

import (
	"errors"

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

func (clients Clients) GenerateDatabases(folders []*models.Folder) error {
	// Generate new client databases
	for _, client := range clients {
		if err := client.GenerateCache(folders); err != nil {
			return err
		}
	}

	return nil
}
