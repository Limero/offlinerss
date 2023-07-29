package mock

import (
	"github.com/limero/go-newsblur"
	"github.com/stretchr/testify/mock"
)

type MockNewsblurClient struct {
	mock.Mock
}

func (m *MockNewsblurClient) Login(username, password string) (*newsblur.LoginOutput, error) {
	args := m.Called(username, password)
	return args.Get(0).(*newsblur.LoginOutput), args.Error(1)
}

func (m *MockNewsblurClient) ReaderFeeds() (*newsblur.ReaderFeedsOutput, error) {
	args := m.Called()
	return args.Get(0).(*newsblur.ReaderFeedsOutput), args.Error(1)
}

func (m *MockNewsblurClient) ReaderUnreadStoryHashes() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockNewsblurClient) ReaderStarredStoryHashes() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockNewsblurClient) ReaderRiverStories_StoryHash(storyHash []string) (*newsblur.StoriesOutput, error) {
	args := m.Called(storyHash)
	return args.Get(0).(*newsblur.StoriesOutput), args.Error(1)
}

func (m *MockNewsblurClient) MarkStoryHashesAsRead(storyHash []string) (*newsblur.MarkStoryHashesAsReadOutput, error) {
	args := m.Called(storyHash)
	return args.Get(0).(*newsblur.MarkStoryHashesAsReadOutput), args.Error(1)
}

func (m *MockNewsblurClient) MarkStoryHashAsUnread(storyHash string) (*newsblur.MarkStoryHashAsUnreadOutput, error) {
	args := m.Called(storyHash)
	return args.Get(0).(*newsblur.MarkStoryHashAsUnreadOutput), args.Error(1)
}

func (m *MockNewsblurClient) MarkStoryHashAsStarred(storyHash string) (*newsblur.MarkStoryHashAsStarredOutput, error) {
	args := m.Called(storyHash)
	return args.Get(0).(*newsblur.MarkStoryHashAsStarredOutput), args.Error(1)
}

func (m *MockNewsblurClient) MarkStoryHashAsUnstarred(storyHash string) (*newsblur.MarkStoryHashAsUnstarredOutput, error) {
	args := m.Called(storyHash)
	return args.Get(0).(*newsblur.MarkStoryHashAsUnstarredOutput), args.Error(1)
}
