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

	if syncToActions.Total() > 0 {
		log.Info("Syncing changes to " + string(server.Name()))

		// TODO: Do these in parallel
		if len(syncToActions.Read) > 0 {
			log.Debug("Syncing read to " + string(server.Name()))
			if err := server.MarkStoriesAsRead(syncToActions.Read); err != nil {
				return nil, err
			}
		}
		if len(syncToActions.Read) > 0 {
			log.Debug("Syncing unread to " + string(server.Name()))
			if err := server.MarkStoriesAsUnread(syncToActions.Unread); err != nil {
				return nil, err
			}
		}
		if len(syncToActions.Starred) > 0 {
			log.Debug("Syncing starred to " + string(server.Name()))
			if err := server.MarkStoriesAsStarred(syncToActions.Starred); err != nil {
				return nil, err
			}
		}
		if len(syncToActions.Unstarred) > 0 {
			log.Debug("Syncing unstarred to " + string(server.Name()))
			if err := server.MarkStoriesAsUnstarred(syncToActions.Unstarred); err != nil {
				return nil, err
			}
		}
	}

	log.Info("Retrieving new stories from " + string(server.Name()))
	return server.GetFoldersWithStories()
}
