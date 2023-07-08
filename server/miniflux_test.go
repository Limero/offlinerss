package server

import (
	"testing"

	"github.com/limero/offlinerss/models"
	"github.com/stretchr/testify/require"
	miniflux "miniflux.app/client"
)

func TestMinifluxGetFoldersWithStories(t *testing.T) {
	t.Skip("TODO")
}

func TestMinifluxSyncToServer(t *testing.T) {
	mockClient := new(MockMinifluxClient)

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
