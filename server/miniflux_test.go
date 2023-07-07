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

	mockClient.On("UpdateEntries", []int64{1, 2}, miniflux.EntryStatusRead).
		Return(nil)
	mockClient.On("UpdateEntries", []int64{3, 4}, miniflux.EntryStatusUnread).
		Return(nil)

	syncToActions := models.SyncToActions{
		{Id: "1", Action: models.ActionStoryRead},
		{Id: "2", Action: models.ActionStoryRead},
		{Id: "3", Action: models.ActionStoryUnread},
		{Id: "4", Action: models.ActionStoryUnread},
	}
	err := s.SyncToServer(syncToActions)
	require.NoError(t, err)

	mockClient.AssertExpectations(t)
}
