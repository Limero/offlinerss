package main

import (
	"errors"
	"os"

	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
)

func GetSyncToActions(clients models.Clients) ([]models.SyncToAction, error) {
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

		if len(actions) == 0 {
			continue
		}

		readActions := 0
		unreadActions := 0
		starredActions := 0
		unstarredActions := 0

		for _, action := range actions {
			switch action.Action {
			case models.ActionStoryRead:
				readActions++
			case models.ActionStoryUnread:
				unreadActions++
			case models.ActionStoryStarred:
				starredActions++
			case models.ActionStoryUnstarred:
				unstarredActions++
			}
		}

		log.Info("Found %d changes in %s", len(actions), client.Name())
		if readActions > 0 {
			log.Info("  📖 %d read", readActions)
		}
		if unreadActions > 0 {
			log.Info("  📕 %d unread", unreadActions)
		}
		if starredActions > 0 {
			log.Info("  ⭐ %d starred", starredActions)
		}
		if unstarredActions > 0 {
			log.Info("  ☁️  %d unstarred", unstarredActions)
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
