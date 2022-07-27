package main

import (
	"github.com/limero/offlinerss/models"
)

type Client interface {
	GetChanges() ([]models.SyncToAction, error)
	GenerateCache(folders []*models.Folder) error
}

type Clients []Client

type Server interface {
	Name() string
	Login() error
	GetFoldersWithStories() ([]*models.Folder, error)
	SyncToServer(syncToActions []models.SyncToAction) error
}
