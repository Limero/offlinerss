package main

import "github.com/limero/offlinerss/models"

type Client interface {
	GetChanges() ([]models.SyncToAction, error)
	GenerateCache(folders []*models.Folder) error
}

type Clients []Client
