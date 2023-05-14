package main

import (
	"github.com/limero/offlinerss/models"
)

type Client interface {
	Name() string
	GetChanges() ([]models.SyncToAction, error)
	CreateNewCache() error
	AddToCache(folders []*models.Folder) error
}

type Clients []Client

type Server interface {
	Name() string
	Login() error
	GetFoldersWithStories() ([]*models.Folder, error)
	SyncToServer(syncToActions []models.SyncToAction) error
}
