package server

import (
	"fmt"
	"testing"
	"time"

	"github.com/limero/offlinerss/models"
	"github.com/limero/offlinerss/server/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	miniflux "miniflux.app/client"
)

func TestMinifluxGetFoldersWithStories(t *testing.T) {
	mockClient := new(mock.MockMinifluxClient)

	s := Miniflux{
		client: mockClient,
	}

	now := time.Now()
	story := models.Story{
		Timestamp: fmt.Sprintf("%d", now.Unix()),
		Hash:      "123",
		Unread:    true,
		Date:      now.Format("2006-01-02 15:04:05"),
	}

	entries := miniflux.Entries{
		{
			ID:     123,
			Status: miniflux.EntryStatusUnread,
			Date:   now,
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
	mockClient := new(mock.MockMinifluxClient)

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
		{Id: "1", Action: models.ActionStoryRead},
		{Id: "2", Action: models.ActionStoryRead},

		{Id: "3", Action: models.ActionStoryUnread},
		{Id: "4", Action: models.ActionStoryUnread},

		{Id: "1", Action: models.ActionStoryStarred},
		{Id: "2", Action: models.ActionStoryStarred}, // already starred, should be skipped

		{Id: "3", Action: models.ActionStoryUnstarred},
		{Id: "4", Action: models.ActionStoryUnstarred}, // already unstarred, should be skipped
	}
	err := s.SyncToServer(syncToActions)
	require.NoError(t, err)

	mockClient.AssertExpectations(t)
}
