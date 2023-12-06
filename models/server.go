package models

type ServerName string

const (
	ServerMiniflux ServerName = "miniflux"
	ServerNewsBlur ServerName = "newsblur"
)

type Server interface {
	Name() ServerName
	Login() error
	GetFoldersWithStories() (Folders, error)
	SyncToServer(syncToActions SyncToActions) error
}
