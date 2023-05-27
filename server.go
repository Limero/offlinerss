package main

import (
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
)

func SyncServer(server models.Server, syncToActions models.SyncToActions) ([]*models.Folder, error) {
	// Sync changes back to server and get new stories

	log.Debug("Logging in to " + server.Name())
	if err := server.Login(); err != nil {
		return nil, err
	}

	if len(syncToActions) > 0 {
		log.Info("Syncing changes to " + server.Name())
		if err := server.SyncToServer(syncToActions); err != nil {
			return nil, err
		}
	}

	log.Info("Retrieving new stories from " + server.Name())
	return server.GetFoldersWithStories()
}
