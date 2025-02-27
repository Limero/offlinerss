package miniflux

import (
	"testing"
	"time"

	"github.com/limero/offlinerss/models"
	"github.com/limero/offlinerss/server/miniflux/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "miniflux.app/v2/client"
)

func TestGetFoldersWithStories(t *testing.T) {
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

func TestMarkStoriesAsRead(t *testing.T) {
	mockAPI := new(mock.MockAPI)

	s := Miniflux{
		api: mockAPI,
	}

	mockAPI.On("UpdateEntries", []int64{1, 2}, api.EntryStatusRead).
		Return(nil)

	require.NoError(t, s.MarkStoriesAsRead([]string{"1", "2"}))
	mockAPI.AssertExpectations(t)
}

func TestMarkStoriesAsUnread(t *testing.T) {
	mockAPI := new(mock.MockAPI)

	s := Miniflux{
		api: mockAPI,
	}

	mockAPI.On("UpdateEntries", []int64{1, 2}, api.EntryStatusUnread).
		Return(nil)

	require.NoError(t, s.MarkStoriesAsUnread([]string{"1", "2"}))
	mockAPI.AssertExpectations(t)
}

func TestMarkStoriesAsStarred(t *testing.T) {
	mockAPI := new(mock.MockAPI)

	s := Miniflux{
		api: mockAPI,
	}

	mockAPI.On("Entry", int64(1)).
		Return(&api.Entry{Starred: false}, nil)
	mockAPI.On("ToggleBookmark", int64(1)).
		Return(nil)
	mockAPI.On("Entry", int64(2)).
		Return(&api.Entry{Starred: true}, nil)

	require.NoError(t, s.MarkStoriesAsStarred([]string{"1", "2"}))
	mockAPI.AssertExpectations(t)
}

func TestMarkStoriesAsUnstarred(t *testing.T) {
	mockAPI := new(mock.MockAPI)

	s := Miniflux{
		api: mockAPI,
	}

	mockAPI.On("Entry", int64(1)).
		Return(&api.Entry{Starred: true}, nil)
	mockAPI.On("ToggleBookmark", int64(1)).
		Return(nil)
	mockAPI.On("Entry", int64(2)).
		Return(&api.Entry{Starred: false}, nil)

	require.NoError(t, s.MarkStoriesAsUnstarred([]string{"1", "2"}))
	mockAPI.AssertExpectations(t)
}
