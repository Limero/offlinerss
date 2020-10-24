package main

import (
	"errors"
	"fmt"

	"github.com/limero/offlinerss/models"
	"github.com/limero/offlinerss/servers/miniflux"
	"github.com/limero/offlinerss/servers/newsblur"
)

func DoSync(serverConfig models.ServerConfig, syncToActions []models.SyncToAction) ([]*models.Folder, error) {
	// Sync changes back to server and get new stories
	var folders []*models.Folder
	switch serverConfig.Type {
	case "miniflux":
		fmt.Println("Logging in to Miniflux")
		minifluxClient, err := miniflux.Login(serverConfig.Username, serverConfig.Password)
		if err != nil {
			return nil, err
		}

		fmt.Println("Syncing changes to Miniflux")
		if err := miniflux.SyncToServer(minifluxClient, syncToActions); err != nil {
			return nil, err
		}

		fmt.Println("Retrieving new stories from Miniflux")
		folders, err = miniflux.GetFoldersWithStories(minifluxClient)
		if err != nil {
			return nil, err
		}
	case "newsblur":
		fmt.Println("Logging in to NewsBlur")
		newsBlurClient, err := newsblur.Login(serverConfig.Username, serverConfig.Password)
		if err != nil {
			return nil, err
		}

		fmt.Println("Syncing changes to NewsBlur")
		if err := newsblur.SyncToServer(newsBlurClient, syncToActions); err != nil {
			return nil, err
		}

		fmt.Println("Retrieving new stories from NewsBlur")
		folders, err = newsblur.GetFoldersWithStories(newsBlurClient)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprintf("Invalid server type: %s", serverConfig.Type))
	}

	return folders, nil
}
