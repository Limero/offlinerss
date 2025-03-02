package domain

import "time"

type ServerName string

const (
	ServerMiniflux ServerName = "miniflux"
	ServerNewsBlur ServerName = "newsblur"
)

type Server interface {
	Name() ServerName
	Login() error
	GetFoldersWithStories(from *time.Time) (Folders, error)

	MarkStoriesAsRead(IDs []string) error
	MarkStoriesAsUnread(IDs []string) error
	MarkStoriesAsStarred(IDs []string) error
	MarkStoriesAsUnstarred(IDs []string) error
}
