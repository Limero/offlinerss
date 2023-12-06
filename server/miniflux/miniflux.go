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
	var readIds []int64
	var unreadIds []int64

	for _, syncToAction := range syncToActions {
		actionId, err := strconv.ParseInt(syncToAction.Id, 10, 64)
		if err != nil {
			return err
		}

		switch syncToAction.Action {
		case models.ActionStoryRead:
			// Batch read events so only one request has to be done
			readIds = append(readIds, actionId)
		case models.ActionStoryUnread:
			// Batch unread events so only one request has to be done
			unreadIds = append(unreadIds, actionId)
		case models.ActionStoryStarred, models.ActionStoryUnstarred:
			if err := s.handleStarred(syncToAction); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unsupported Miniflux syncToAction: %d", syncToAction.Action)
		}
	}

	if len(readIds) > 0 {
		if err := s.client.UpdateEntries(readIds, miniflux.EntryStatusRead); err != nil {
			return err
		}
		log.Debug("%d items has been marked as read", len(readIds))
	}

	if len(unreadIds) > 0 {
		if err := s.client.UpdateEntries(unreadIds, miniflux.EntryStatusUnread); err != nil {
			return err
		}
		log.Debug("%d items has been marked as unread", len(unreadIds))
	}

	return nil
}

func (s *Miniflux) handleStarred(syncToAction models.SyncToAction) error {
	// Because Miniflux only support toggling starred instead of setting it directly,
	// we have to check its current status

	actionId, err := strconv.ParseInt(syncToAction.Id, 10, 64)
	if err != nil {
		return err
	}

	entry, err := s.client.Entry(actionId)
	if err != nil {
		return err
	}

	if (entry.Starred && syncToAction.Action == models.ActionStoryUnstarred) ||
		(!entry.Starred && syncToAction.Action == models.ActionStoryStarred) {
		return s.client.ToggleBookmark(actionId)
	}

	return nil
}