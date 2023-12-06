package miniflux

import (
	"testing"
	"time"

	"github.com/limero/offlinerss/models"
	"github.com/limero/offlinerss/server/miniflux/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	miniflux "miniflux.app/client"
)

func TestMinifluxGetFoldersWithStories(t *testing.T) {
	mockClient := new(mock.MockClient)

	s := Miniflux{
		client: mockClient,
	}

	story := models.Story{
		Timestamp: time.Now(),
		Hash:      "123",
		Unread:    true,
	}

	entries := miniflux.Entries{
		{
			ID:     123,
			Status: miniflux.EntryStatusUnread,
			Date:   story.Timestamp,
			Feed: &miniflux.Feed{
				Category: &miniflux.Category{},
			},
		},
	}

	mockClient.On("Entries", &miniflux.Filter{
		Status: miniflux.EntryStatusUnread,
	}).Return(&miniflux.EntryResultSet{
		Total:   len(entries),
		Entries: entries,
	}, nil)

	folders, err := s.GetFoldersWithStories()
	require.NoError(t, err)

	assert.Len(t, folders, 1)
	assert.Len(t, folders[0].Feeds, 1)
	assert.Len(t, folders[0].Feeds[0].Stories, 1)
	assert.Equal(t, &story, folders[0].Feeds[0].Stories[0])

	mockClient.AssertExpectations(t)
}

func TestMinifluxSyncToServer(t *testing.T) {
	mockClient := new(mock.MockClient)

	s := Miniflux{
		client: mockClient,
	}

	// Read
	mockClient.On("UpdateEntries", []int64{1, 2}, miniflux.EntryStatusRead).
		Return(nil)

	// Unread
	mockClient.On("UpdateEntries", []int64{3, 4}, miniflux.EntryStatusUnread).
		Return(nil)

	// Starred
	mockClient.On("Entry", int64(1)).
		Return(&miniflux.Entry{Starred: false}, nil)
	mockClient.On("ToggleBookmark", int64(1)).
		Return(nil)
	mockClient.On("Entry", int64(2)).
		Return(&miniflux.Entry{Starred: true}, nil)

	// Unstarred
	mockClient.On("Entry", int64(3)).
		Return(&miniflux.Entry{Starred: true}, nil)
	mockClient.On("ToggleBookmark", int64(3)).
		Return(nil)
	mockClient.On("Entry", int64(4)).
		Return(&miniflux.Entry{Starred: false}, nil)

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

	mockClient.AssertExpectations(t)
}
