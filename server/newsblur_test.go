package server

import (
	"testing"
	"time"

	"github.com/limero/go-newsblur"
	"github.com/limero/offlinerss/models"
	"github.com/limero/offlinerss/server/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewsblurGetFoldersWithStories(t *testing.T) {
	mockClient := new(mock.MockNewsblurClient)

	s := Newsblur{
		client: mockClient,
	}
	now := time.Now()
	story := models.Story{
		Timestamp: time.Unix(now.Unix(), 0),
		Hash:      "123",
		Unread:    true,
	}

	mockClient.On("ReaderFeeds").
		Return(&newsblur.ReaderFeedsOutput{
			Folders: []newsblur.Folder{
				{
					Title:   "folder",
					FeedIDs: []int{1},
				},
			},
			Feeds: []newsblur.ApiFeed{
				{
					ID: 1,
					Nt: 1,
				},
			},
		}, nil)

	mockClient.On("ReaderRiverStories", []string{"1"}, 1).
		Return(&newsblur.StoriesOutput{
			Stories: []newsblur.ApiStory{
				{
					StoryTimestamp: story.Timestamp.Unix(),
					StoryHash:      story.Hash,
					StoryFeedID:    1,
				},
			},
		}, nil)

	mockClient.On("ReaderStarredStories", 1).
		Return(&newsblur.StoriesOutput{
			Stories: []newsblur.ApiStory{
				{
					StoryTimestamp: now.Unix(),
					StoryHash:      "321",
					StoryFeedID:    1,
				},
			},
		}, nil)
	mockClient.On("ReaderStarredStories", 2).
		Return(&newsblur.StoriesOutput{
			Stories: []newsblur.ApiStory{},
		}, nil)

	folders, err := s.GetFoldersWithStories()
	require.NoError(t, err)

	assert.Len(t, folders, 1)
	assert.Len(t, folders[0].Feeds, 1)
	assert.Len(t, folders[0].Feeds[0].Stories, 2)
	assert.Equal(t, &story, folders[0].Feeds[0].Stories[0])

	mockClient.AssertExpectations(t)
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
