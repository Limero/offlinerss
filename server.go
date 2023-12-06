package main

import (
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
	"github.com/limero/offlinerss/server/miniflux"
	"github.com/limero/offlinerss/server/newsblur"
)

func getServer(serverConfig models.ServerConfig) models.Server {
	switch serverConfig.Name {
	case models.ServerMiniflux:
		return miniflux.New(serverConfig)
	case models.ServerNewsBlur:
		return newsblur.New(serverConfig)
	}
	return nil
}

func SyncServer(server models.Server, syncToActions models.SyncToActions) (models.Folders, error) {
	// Sync changes back to server and get new stories

	log.Debug("Logging in to " + string(server.Name()))
	if err := server.Login(); err != nil {
		return nil, err
	}

	if len(syncToActions) > 0 {
		log.Info("Syncing changes to " + string(server.Name()))
		if err := server.SyncToServer(syncToActions); err != nil {
			return nil, err
		}
	}

	log.Info("Retrieving new stories from " + string(server.Name()))
	return server.GetFoldersWithStories()
}
