package main

import (
	"github.com/limero/offlinerss/domain"
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/server/miniflux"
	"github.com/limero/offlinerss/server/newsblur"
)

func getServer(serverConfig domain.ServerConfig) domain.Server {
	switch serverConfig.Name {
	case domain.ServerMiniflux:
		return miniflux.New(serverConfig)
	case domain.ServerNewsBlur:
		return newsblur.New(serverConfig)
	}
	return nil
}

func AuthServer(server domain.Server) error {
	log.Debug("Logging in to " + string(server.Name()))
	return server.Login()
}

func SyncToServer(server domain.Server, syncToActions domain.SyncToActions) error {
	log.Info("Syncing changes to " + string(server.Name()))

	// TODO: Do these in parallel
	if len(syncToActions.Read) > 0 {
		log.Debug("Syncing read to " + string(server.Name()))
		if err := server.MarkStoriesAsRead(syncToActions.Read); err != nil {
			return err
		}
	}
	if len(syncToActions.Unread) > 0 {
		log.Debug("Syncing unread to " + string(server.Name()))
		if err := server.MarkStoriesAsUnread(syncToActions.Unread); err != nil {
			return err
		}
	}
	if len(syncToActions.Starred) > 0 {
		log.Debug("Syncing starred to " + string(server.Name()))
		if err := server.MarkStoriesAsStarred(syncToActions.Starred); err != nil {
			return err
		}
	}
	if len(syncToActions.Unstarred) > 0 {
		log.Debug("Syncing unstarred to " + string(server.Name()))
		if err := server.MarkStoriesAsUnstarred(syncToActions.Unstarred); err != nil {
			return err
		}
	}

	return nil
}
