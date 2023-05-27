package models

type Server interface {
	Name() string
	Login() error
	GetFoldersWithStories() (Folders, error)
	SyncToServer(syncToActions SyncToActions) error
}
