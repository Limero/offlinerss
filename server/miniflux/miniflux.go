package miniflux

import (
	"fmt"
	"strconv"

	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
	miniflux "miniflux.app/client"
)

type MinifluxClient interface {
	Entry(entryID int64) (*miniflux.Entry, error)
	Entries(filter *miniflux.Filter) (*miniflux.EntryResultSet, error)
	UpdateEntries(entryIDs []int64, status string) error
	ToggleBookmark(entryID int64) error
}

type Miniflux struct {
	config models.ServerConfig
	client MinifluxClient
}

func New(config models.ServerConfig) *Miniflux {
	return &Miniflux{
		config: config,
	}
}

func (s *Miniflux) Name() string {
	return s.config.Type
}

func (s *Miniflux) Login() error {
	client := miniflux.New("http://localhost", s.config.Username, s.config.Password)

	if _, err := client.Me(); err != nil {
		return nil
	}

	s.client = client
	return nil
}

func (s *Miniflux) GetFoldersWithStories() (models.Folders, error) {
	var folders models.Folders

	entries, err := s.client.Entries(&miniflux.Filter{
		Status: miniflux.EntryStatusUnread,
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
			Unread:    entry.Status != miniflux.EntryStatusRead,
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

func (s *Miniflux) SyncToServer(syncToActions models.SyncToActions) error {
	var readIDs []int64
	var unreadIDs []int64

	for _, syncToAction := range syncToActions {
		actionID, err := strconv.ParseInt(syncToAction.ID, 10, 64)
		if err != nil {
			return err
		}

		switch syncToAction.Action {
		case models.ActionStoryRead:
			// Batch read events so only one request has to be done
			readIDs = append(readIDs, actionID)
		case models.ActionStoryUnread:
			// Batch unread events so only one request has to be done
			unreadIDs = append(unreadIDs, actionID)
		case models.ActionStoryStarred, models.ActionStoryUnstarred:
			if err := s.handleStarred(syncToAction); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unsupported Miniflux syncToAction: %d", syncToAction.Action)
		}
	}

	if len(readIDs) > 0 {
		if err := s.client.UpdateEntries(readIDs, miniflux.EntryStatusRead); err != nil {
			return err
		}
		log.Debug("%d items has been marked as read", len(readIDs))
	}

	if len(unreadIDs) > 0 {
		if err := s.client.UpdateEntries(unreadIDs, miniflux.EntryStatusUnread); err != nil {
			return err
		}
		log.Debug("%d items has been marked as unread", len(unreadIDs))
	}

	return nil
}

func (s *Miniflux) handleStarred(syncToAction models.SyncToAction) error {
	// Because Miniflux only support toggling starred instead of setting it directly,
	// we have to check its current status

	actionID, err := strconv.ParseInt(syncToAction.ID, 10, 64)
	if err != nil {
		return err
	}

	entry, err := s.client.Entry(actionID)
	if err != nil {
		return err
	}

	if (entry.Starred && syncToAction.Action == models.ActionStoryUnstarred) ||
		(!entry.Starred && syncToAction.Action == models.ActionStoryStarred) {
		return s.client.ToggleBookmark(actionID)
	}

	return nil
}
