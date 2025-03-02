package domain

type ServerName string

const (
	ServerMiniflux ServerName = "miniflux"
	ServerNewsBlur ServerName = "newsblur"
)

type Server interface {
	Name() ServerName
	Login() error
	GetFoldersWithStories() (Folders, error)

	MarkStoriesAsRead(IDs []string) error
	MarkStoriesAsUnread(IDs []string) error
	MarkStoriesAsStarred(IDs []string) error
	MarkStoriesAsUnstarred(IDs []string) error
}
