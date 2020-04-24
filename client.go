package main

import (
	"errors"
	"fmt"

	"github.com/limero/offlinerss/clients/feedreader"
	"github.com/limero/offlinerss/clients/newsboat"
	"github.com/limero/offlinerss/models"
)

func GetSyncToActions(clientConfigs []models.ClientConfig) ([]models.SyncToAction, error) {
	// Grab changes from each client
	var syncToActions []models.SyncToAction
	hasEnabledClient := false

	for _, clientConfig := range clientConfigs {
		if !clientConfig.Enabled {
			continue
		}
		switch clientConfig.Type {
		case "newsboat":
			actions, err := newsboat.GetChanges(clientConfig)
			if err != nil {
				return nil, err
			}
			syncToActions = append(syncToActions, actions...)
		case "feedreader":
			actions, err := feedreader.GetChanges(clientConfig)
			if err != nil {
				return nil, err
			}
			syncToActions = append(syncToActions, actions...)
		default:
			return nil, errors.New(fmt.Sprintf("Invalid client type: %s", clientConfig.Type))
		}

		hasEnabledClient = true
	}

	if !hasEnabledClient {
		return nil, errors.New("You have to enable at least one client in the config file")
	}

	return syncToActions, nil
}

func GenerateDatabases(clientConfigs []models.ClientConfig, folders []*models.Folder) error {
	// Generate new client databases
	for _, clientConfig := range clientConfigs {
		if !clientConfig.Enabled {
			continue
		}
		switch clientConfig.Type {
		case "newsboat":
			if err := newsboat.GenerateCache(folders, clientConfig); err != nil {
				return err
			}
		case "feedreader":
			if err := feedreader.GenerateCache(folders, clientConfig); err != nil {
				return err
			}
		default:
			return errors.New(fmt.Sprintf("Invalid client type: %s", clientConfig.Type))
		}
	}

	return nil
}
