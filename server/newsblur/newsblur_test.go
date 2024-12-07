package newsblur

import (
	"testing"
	"time"

	"github.com/limero/offlinerss/models"
	newsblur "github.com/limero/offlinerss/server/newsblur/api"
	"github.com/limero/offlinerss/server/newsblur/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewsblurGetFoldersWithStories(t *testing.T) {
	mockClient := new(mock.MockClient)

	s := Newsblur{
		client: mockClient,
	}
	now := time.Now()
	unreadStory := models.Story{
		Timestamp: time.Unix(now.Unix(), 0),
		Hash:      "123",
		Unread:    true,
	}
	starredStory := models.Story{
		Timestamp: time.Unix(now.Unix(), 0),
		Hash:      "321",
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

	mockClient.On("ReaderUnreadStoryHashes").
		Return([]string{unreadStory.Hash}, nil)

	mockClient.On("ReaderStarredStoryHashes").
		Return([]string{starredStory.Hash}, nil)

	mockClient.On("ReaderRiverStories_StoryHash", []string{unreadStory.Hash, starredStory.Hash}).
		Return(&newsblur.StoriesOutput{
			Stories: []newsblur.ApiStory{
				{
					StoryTimestamp: unreadStory.Timestamp.Unix(),
					StoryHash:      unreadStory.Hash,
					StoryFeedID:    1,
				},
				{
					StoryTimestamp: starredStory.Timestamp.Unix(),
					StoryHash:      starredStory.Hash,
					StoryFeedID:    1,
				},
			},
		}, nil)

	folders, err := s.GetFoldersWithStories()
	require.NoError(t, err)

	assert.Len(t, folders, 1)
	assert.Len(t, folders[0].Feeds, 1)
	assert.Len(t, folders[0].Feeds[0].Stories, 2)
	assert.Equal(t, &unreadStory, folders[0].Feeds[0].Stories[0])
	assert.Equal(t, &starredStory, folders[0].Feeds[0].Stories[1])

	mockClient.AssertExpectations(t)
}

func TestNewsblurSyncToServer(t *testing.T) {
	mockClient := new(mock.MockClient)

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
		{ID: "1", Action: models.ActionStoryRead},
		{ID: "2", Action: models.ActionStoryRead},

		{ID: "3", Action: models.ActionStoryUnread},
		{ID: "4", Action: models.ActionStoryUnread},

		{ID: "1", Action: models.ActionStoryStarred},
		{ID: "2", Action: models.ActionStoryStarred},

		{ID: "3", Action: models.ActionStoryUnstarred},
		{ID: "4", Action: models.ActionStoryUnstarred},
	}
	err := s.SyncToServer(syncToActions)
	require.NoError(t, err)

	mockClient.AssertExpectations(t)
}
