package models

type Server interface {
	Name() string
	Login() error
	GetFoldersWithStories() ([]*Folder, error)
	SyncToServer(syncToActions SyncToActions) error
}
