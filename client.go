package main

import (
	"errors"

	"github.com/limero/offlinerss/client/feedreader"
	"github.com/limero/offlinerss/client/newsboat"
	"github.com/limero/offlinerss/client/quiterss"
	"github.com/limero/offlinerss/domain"
	"github.com/limero/offlinerss/helpers"
	"github.com/limero/offlinerss/log"
)

func getClients(clientConfigs []domain.ClientConfig) domain.Clients {
	clients := make(domain.Clients, len(clientConfigs))
	for i, clientConfig := range clientConfigs {
		switch clientConfig.Name {
		case domain.ClientFeedReader:
			clients[i] = feedreader.New(clientConfig)
		case domain.ClientNewsboat:
			clients[i] = newsboat.New(clientConfig)
		case domain.ClientQuiteRSS:
			clients[i] = quiterss.New(clientConfig)
		default:
			panic("unknown client " + clientConfig.Name)
		}
	}
	return clients
}

func GetSyncToActions(clients domain.Clients) (domain.SyncToActions, error) {
	if len(clients) == 0 {
		return domain.SyncToActions{}, errors.New("you have to enable at least one client in the config file")
	}

	// Grab changes from each client
	var syncToActions domain.SyncToActions

	for _, client := range clients {
		actions, err := client.GetChanges()
		if err != nil {
			return domain.SyncToActions{}, err
		}

		if actions.Total() == 0 {
			continue
		}

		log.Info("Found %d changes in %s", actions.Total(), client.Name())
		if len(actions.Read) > 0 {
			log.Info("  📖 %d read", len(actions.Read))
			syncToActions.Read = append(syncToActions.Read, actions.Read...)
		}
		if len(actions.Unread) > 0 {
			log.Info("  📕 %d unread", len(actions.Unread))
			syncToActions.Unread = append(syncToActions.Unread, actions.Unread...)
		}
		if len(actions.Starred) > 0 {
			log.Info("  ⭐ %d starred", len(actions.Starred))
			syncToActions.Starred = append(syncToActions.Starred, actions.Starred...)
		}
		if len(actions.Unstarred) > 0 {
			log.Info("  ☁️ %d unstarred", len(actions.Unstarred))
			syncToActions.Unstarred = append(syncToActions.Unstarred, actions.Unstarred...)
		}
	}

	if syncToActions.Total() == 0 {
		log.Info("No changes found for any local clients")
	}

	return syncToActions, nil
}

func SyncClients(clients domain.Clients, folders domain.Folders) error {
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

// This will replace all reference dbs with user dbs
// Needed when only syncing to a server, since otherwise the same
// changes will be synced again on next run
func ReplaceReferenceDBsWithUserDBs(clients domain.Clients) error {
	for _, client := range clients {
		log.Debug("Replacing %s reference db with user db", client.Name())
		if err := helpers.CopyFile(client.UserDB(), client.ReferenceDB()); err != nil {
			return err
		}
	}
	return nil
}

// This will replace all user dbs with reference dbs
// It's used to rollback any changes done to the clients since the last sync
func ReplaceUserDBsWithReferenceDBs(clients domain.Clients) error {
	for _, client := range clients {
		log.Debug("Replacing %s user db with reference db", client.Name())
		if err := helpers.CopyFile(client.ReferenceDB(), client.UserDB()); err != nil {
			return err
		}
	}
	return nil
}
