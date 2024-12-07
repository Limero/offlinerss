package mock

import (
	newsblur "github.com/limero/offlinerss/server/newsblur/api"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Login(username, password string) (*newsblur.LoginOutput, error) {
	args := m.Called(username, password)
	return args.Get(0).(*newsblur.LoginOutput), args.Error(1)
}

func (m *MockClient) ReaderFeeds() (*newsblur.ReaderFeedsOutput, error) {
	args := m.Called()
	return args.Get(0).(*newsblur.ReaderFeedsOutput), args.Error(1)
}

func (m *MockClient) ReaderUnreadStoryHashes() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockClient) ReaderStarredStoryHashes() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockClient) ReaderRiverStories_StoryHash(storyHash []string) (*newsblur.StoriesOutput, error) {
	args := m.Called(storyHash)
	return args.Get(0).(*newsblur.StoriesOutput), args.Error(1)
}

func (m *MockClient) MarkStoryHashesAsRead(storyHash []string) (*newsblur.MarkStoryHashesAsReadOutput, error) {
	args := m.Called(storyHash)
	return args.Get(0).(*newsblur.MarkStoryHashesAsReadOutput), args.Error(1)
}

func (m *MockClient) MarkStoryHashAsUnread(storyHash string) (*newsblur.MarkStoryHashAsUnreadOutput, error) {
	args := m.Called(storyHash)
	return args.Get(0).(*newsblur.MarkStoryHashAsUnreadOutput), args.Error(1)
}

func (m *MockClient) MarkStoryHashAsStarred(storyHash string) (*newsblur.MarkStoryHashAsStarredOutput, error) {
	args := m.Called(storyHash)
	return args.Get(0).(*newsblur.MarkStoryHashAsStarredOutput), args.Error(1)
}

func (m *MockClient) MarkStoryHashAsUnstarred(storyHash string) (*newsblur.MarkStoryHashAsUnstarredOutput, error) {
	args := m.Called(storyHash)
	return args.Get(0).(*newsblur.MarkStoryHashAsUnstarredOutput), args.Error(1)
}
