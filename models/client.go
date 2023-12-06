package models

type ClientName string

const (
	ClientFeedReader ClientName = "feedreader"
	ClientNewsboat   ClientName = "newsboat"
	ClientQuiteRSS   ClientName = "quiterss"
)

type Client interface {
	Name() ClientName
	UserDB() string
	ReferenceDB() string
	GetChanges() (SyncToActions, error)
	GetDatabaseInfo() DatabaseInfo
	CreateNewCache() error
	AddToCache(folders Folders) error
	SetDataPath(dataPath DataPath)
}

type Clients []Client
