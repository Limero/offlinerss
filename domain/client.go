package domain

type ClientName string

const (
	ClientFeedReader ClientName = "feedreader"
	ClientNewsboat   ClientName = "newsboat"
	ClientNewsraft   ClientName = "newsraft"
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
	GetDataPath() DataPath
	SetDataPath(dataPath DataPath)
	GetFiles() ClientFiles
}

type Clients []Client

type ClientFile struct {
	FileName    string
	TargetPaths []string
}

type ClientFiles []ClientFile
