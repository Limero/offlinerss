package newsblur

import (
	"testing"
	"time"

	"github.com/limero/offlinerss/domain"
	"github.com/limero/offlinerss/server/newsblur/api"
	newsblur "github.com/limero/offlinerss/server/newsblur/api"
	"github.com/limero/offlinerss/server/newsblur/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFoldersWithStories(t *testing.T) {
	mockAPI := new(mock.MockAPI)

	s := Newsblur{
		api: mockAPI,
	}
	now := time.Now()
	unreadStory := domain.Story{
		Timestamp: time.Unix(now.Unix(), 0),
		Hash:      "123",
		Unread:    true,
	}
	starredStory := domain.Story{
		Timestamp: time.Unix(now.Unix(), 0),
		Hash:      "321",
		Unread:    true,
	}

	mockAPI.On("ReaderFeeds").
		Return(&newsblur.ReaderFeedsOutput{
			Folders: []newsblur.Folder{
				{
					Title:   "folder",
					FeedIDs: []int{1},
				},
			},
			Feeds: []newsblur.Feed{
				{
					ID: 1,
					Nt: 1,
				},
			},
		}, nil)

	mockAPI.On("ReaderUnreadStoryHashes").
		Return([]string{unreadStory.Hash}, nil)

	mockAPI.On("ReaderStarredStoryHashes").
		Return([]api.HashWithTimestamp{{Hash: starredStory.Hash}}, nil)

	mockAPI.On("ReaderRiverStories_StoryHash", []string{unreadStory.Hash, starredStory.Hash}).
		Return([]newsblur.Story{
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
		}, nil)

	folders, err := s.GetFoldersWithStories(nil)
	require.NoError(t, err)

	assert.Len(t, folders, 1)
	assert.Len(t, folders[0].Feeds, 1)
	assert.Len(t, folders[0].Feeds[0].Stories, 2)
	assert.Equal(t, &unreadStory, folders[0].Feeds[0].Stories[0])
	assert.Equal(t, &starredStory, folders[0].Feeds[0].Stories[1])

	mockAPI.AssertExpectations(t)
}

func TestMarkStoriesAsRead(t *testing.T) {
	mockAPI := new(mock.MockAPI)

	s := Newsblur{
		api: mockAPI,
	}

	mockAPI.On("MarkStoryHashesAsRead", []string{"1", "2"}).
		Return(nil)

	require.NoError(t, s.MarkStoriesAsRead([]string{"1", "2"}))
	mockAPI.AssertExpectations(t)
}

func TestMarkStoriesAsUnread(t *testing.T) {
	mockAPI := new(mock.MockAPI)

	s := Newsblur{
		api: mockAPI,
	}

	mockAPI.On("MarkStoryHashAsUnread", []string{"1", "2"}).
		Return(nil)

	require.NoError(t, s.MarkStoriesAsUnread([]string{"1", "2"}))
	mockAPI.AssertExpectations(t)
}

func TestMarkStoriesAsStarred(t *testing.T) {
	mockAPI := new(mock.MockAPI)

	s := Newsblur{
		api: mockAPI,
	}

	mockAPI.On("MarkStoryHashAsStarred", []string{"1", "2"}).
		Return(nil)

	require.NoError(t, s.MarkStoriesAsStarred([]string{"1", "2"}))
	mockAPI.AssertExpectations(t)
}

func TestMarkStoriesAsUnstarred(t *testing.T) {
	mockAPI := new(mock.MockAPI)

	s := Newsblur{
		api: mockAPI,
	}

	mockAPI.On("MarkStoryHashAsUnstarred", []string{"1", "2"}).
		Return(nil)

	require.NoError(t, s.MarkStoriesAsUnstarred([]string{"1", "2"}))
	mockAPI.AssertExpectations(t)
}
