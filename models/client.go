package models

type Client interface {
	Name() string
	UserDB() string
	ReferenceDB() string
	GetChanges() (SyncToActions, error)
	GetDatabaseInfo() DatabaseInfo
	CreateNewCache() error
	AddToCache(folders Folders) error
}

type Clients []Client
