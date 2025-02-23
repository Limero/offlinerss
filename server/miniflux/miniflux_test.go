package miniflux

import (
	"testing"
	"time"

	"github.com/limero/offlinerss/models"
	"github.com/limero/offlinerss/server/miniflux/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "miniflux.app/client"
)

func TestMinifluxGetFoldersWithStories(t *testing.T) {
	mockAPI := new(mock.MockAPI)

	s := Miniflux{
		api: mockAPI,
	}

	story := models.Story{
		Timestamp: time.Now(),
		Hash:      "123",
		Unread:    true,
	}

	entries := api.Entries{
		{
			ID:     123,
			Status: api.EntryStatusUnread,
			Date:   story.Timestamp,
			Feed: &api.Feed{
				Category: &api.Category{},
			},
		},
	}

	mockAPI.On("Entries", &api.Filter{
		Status: api.EntryStatusUnread,
	}).Return(&api.EntryResultSet{
		Total:   len(entries),
		Entries: entries,
	}, nil)

	folders, err := s.GetFoldersWithStories()
	require.NoError(t, err)

	assert.Len(t, folders, 1)
	assert.Len(t, folders[0].Feeds, 1)
	assert.Len(t, folders[0].Feeds[0].Stories, 1)
	assert.Equal(t, &story, folders[0].Feeds[0].Stories[0])

	mockAPI.AssertExpectations(t)
}

func TestMinifluxSyncToServer(t *testing.T) {
	mockAPI := new(mock.MockAPI)

	s := Miniflux{
		api: mockAPI,
	}

	// Read
	mockAPI.On("UpdateEntries", []int64{1, 2}, api.EntryStatusRead).
		Return(nil)

	// Unread
	mockAPI.On("UpdateEntries", []int64{3, 4}, api.EntryStatusUnread).
		Return(nil)

	// Starred
	mockAPI.On("Entry", int64(1)).
		Return(&api.Entry{Starred: false}, nil)
	mockAPI.On("ToggleBookmark", int64(1)).
		Return(nil)
	mockAPI.On("Entry", int64(2)).
		Return(&api.Entry{Starred: true}, nil)

	// Unstarred
	mockAPI.On("Entry", int64(3)).
		Return(&api.Entry{Starred: true}, nil)
	mockAPI.On("ToggleBookmark", int64(3)).
		Return(nil)
	mockAPI.On("Entry", int64(4)).
		Return(&api.Entry{Starred: false}, nil)

	syncToActions := models.SyncToActions{
		{ID: "1", Action: models.ActionStoryRead},
		{ID: "2", Action: models.ActionStoryRead},

		{ID: "3", Action: models.ActionStoryUnread},
		{ID: "4", Action: models.ActionStoryUnread},

		{ID: "1", Action: models.ActionStoryStarred},
		{ID: "2", Action: models.ActionStoryStarred}, // already starred, should be skipped

		{ID: "3", Action: models.ActionStoryUnstarred},
		{ID: "4", Action: models.ActionStoryUnstarred}, // already unstarred, should be skipped
	}
	err := s.SyncToServer(syncToActions)
	require.NoError(t, err)

	mockAPI.AssertExpectations(t)
}
