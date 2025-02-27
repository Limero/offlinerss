package miniflux

import (
	"strconv"

	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
	api "miniflux.app/v2/client"
)

type API interface {
	Entry(entryID int64) (*api.Entry, error)
	Entries(filter *api.Filter) (*api.EntryResultSet, error)
	UpdateEntries(entryIDs []int64, status string) error
	ToggleBookmark(entryID int64) error
}

type Miniflux struct {
	config models.ServerConfig
	api    API
}

func New(config models.ServerConfig) *Miniflux {
	return &Miniflux{
		config: config,
	}
}

func (s *Miniflux) Name() models.ServerName {
	return s.config.Name
}

func (s *Miniflux) Login() error {
	hostname := s.config.Hostname
	if hostname == "" {
		hostname = "https://reader.miniflux.app"
	}
	client := api.NewClient(hostname, s.config.Username, s.config.Password)

	if _, err := client.Me(); err != nil {
		return err
	}

	s.api = client
	return nil
}

func (s *Miniflux) GetFoldersWithStories() (models.Folders, error) {
	var folders models.Folders

	entries, err := s.api.Entries(&api.Filter{
		Status: api.EntryStatusUnread,
	})
	if err != nil {
		return nil, err
	}

	for _, entry := range entries.Entries {
		story := &models.Story{
			Timestamp: entry.Date,
			Hash:      strconv.FormatInt(entry.ID, 10), // Miniflux has "hash" but IDs are used for marking entries
			Title:     entry.Title,
			Authors:   entry.Author,
			Content:   entry.Content,
			Url:       entry.URL,
			Unread:    entry.Status != api.EntryStatusRead,
			Starred:   entry.Starred,
		}

		storyFolder := folders.GetOrCreateFolder(entry.Feed.Category.ID, entry.Feed.Category.Title)
		storyFeed := storyFolder.Feeds.GetOrCreateFeed(
			entry.Feed.ID,
			entry.Feed.Title,
			entry.Feed.FeedURL,
			entry.Feed.SiteURL,
		)

		if story.Unread {
			storyFeed.Unread++
		}

		storyFeed.Stories = append(storyFeed.Stories, story)
	}

	return folders, nil
}

func (s *Miniflux) MarkStoriesAsRead(IDs []string) error {
	intIDs := make([]int64, 0, len(IDs))

	for _, id := range IDs {
		intID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return err
		}

		intIDs = append(intIDs, intID)
	}

	if err := s.api.UpdateEntries(intIDs, api.EntryStatusRead); err != nil {
		return err
	}
	log.Debug("%d items has been marked as read", len(intIDs))

	return nil
}

func (s *Miniflux) MarkStoriesAsUnread(IDs []string) error {
	intIDs := make([]int64, 0, len(IDs))

	for _, id := range IDs {
		intID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return err
		}

		intIDs = append(intIDs, intID)
	}

	if err := s.api.UpdateEntries(intIDs, api.EntryStatusUnread); err != nil {
		return err
	}
	log.Debug("%d items has been marked as unread", len(intIDs))

	return nil
}

func (s *Miniflux) MarkStoriesAsStarred(IDs []string) error {
	for _, id := range IDs {
		intID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return err
		}

		// Because Miniflux only support toggling starred instead of setting it directly,
		// we have to check its current status
		entry, err := s.api.Entry(intID)
		if err != nil {
			return err
		}
		if entry.Starred {
			continue
		}
		if err := s.api.ToggleBookmark(intID); err != nil {
			return err
		}

	}

	log.Debug("%d items has been marked as starred", len(IDs))

	return nil
}

func (s *Miniflux) MarkStoriesAsUnstarred(IDs []string) error {
	for _, id := range IDs {
		intID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return err
		}

		// Because Miniflux only support toggling starred instead of setting it directly,
		// we have to check its current status
		entry, err := s.api.Entry(intID)
		if err != nil {
			return err
		}
		if !entry.Starred {
			continue
		}
		if err := s.api.ToggleBookmark(intID); err != nil {
			return err
		}

	}

	log.Debug("%d items has been marked as unstarred", len(IDs))

	return nil
}
