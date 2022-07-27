package main

import (
	"fmt"

	"github.com/limero/offlinerss/models"
)

func SyncServer(server Server, syncToActions []models.SyncToAction) ([]*models.Folder, error) {
	// Sync changes back to server and get new stories

	fmt.Println("Logging in to " + server.Name())
	if err := server.Login(); err != nil {
		return nil, err
	}

	fmt.Println("Syncing changes to " + server.Name())
	if err := server.SyncToServer(syncToActions); err != nil {
		return nil, err
	}

	fmt.Println("Retrieving new stories from " + server.Name())
	return server.GetFoldersWithStories()
}
