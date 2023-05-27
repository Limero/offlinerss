package models

type Client interface {
	Name() string
	UserDB() string
	ReferenceDB() string
	GetChanges() (SyncToActions, error)
	CreateNewCache() error
	AddToCache(folders []*Folder) error
}

type Clients []Client
