package models

type Client interface {
	Name() string
	UserDB() string
	ReferenceDB() string
	GetChanges() (SyncToActions, error)
	GetDatabaseInfo() DatabaseInfo
	CreateNewCache() error
	AddToCache(folders Folders) error
	SetDataPath(dataPath DataPath)
}

type Clients []Client
