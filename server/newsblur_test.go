package server

import (
	"testing"

	"github.com/limero/go-newsblur"
	"github.com/limero/offlinerss/models"
	"github.com/limero/offlinerss/server/mock"
	"github.com/stretchr/testify/require"
)

func TestNewsblurGetFoldersWithStories(t *testing.T) {
	t.Skip("TODO")
}

func TestNewsblurSyncToServer(t *testing.T) {
	mockClient := new(mock.MockNewsblurClient)

	s := Newsblur{
		client: mockClient,
	}

	// Read
	mockClient.On("MarkStoryHashesAsRead", []string{"1", "2"}).
		Return(&newsblur.MarkStoryHashesAsReadOutput{}, nil)

	// Unread
	mockClient.On("MarkStoryHashAsUnread", "3").
		Return(&newsblur.MarkStoryHashAsUnreadOutput{}, nil)
	mockClient.On("MarkStoryHashAsUnread", "4").
		Return(&newsblur.MarkStoryHashAsUnreadOutput{}, nil)

	// Starred
	mockClient.On("MarkStoryHashAsStarred", "1").
		Return(&newsblur.MarkStoryHashAsStarredOutput{}, nil)
	mockClient.On("MarkStoryHashAsStarred", "2").
		Return(&newsblur.MarkStoryHashAsStarredOutput{}, nil)

	// Unstarred
	mockClient.On("MarkStoryHashAsUnstarred", "3").
		Return(&newsblur.MarkStoryHashAsUnstarredOutput{}, nil)
	mockClient.On("MarkStoryHashAsUnstarred", "4").
		Return(&newsblur.MarkStoryHashAsUnstarredOutput{}, nil)

	syncToActions := models.SyncToActions{
		{Id: "1", Action: models.ActionStoryRead},
		{Id: "2", Action: models.ActionStoryRead},

		{Id: "3", Action: models.ActionStoryUnread},
		{Id: "4", Action: models.ActionStoryUnread},

		{Id: "1", Action: models.ActionStoryStarred},
		{Id: "2", Action: models.ActionStoryStarred},

		{Id: "3", Action: models.ActionStoryUnstarred},
		{Id: "4", Action: models.ActionStoryUnstarred},
	}
	err := s.SyncToServer(syncToActions)
	require.NoError(t, err)

	mockClient.AssertExpectations(t)
}
