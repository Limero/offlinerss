package mock

import (
	newsblur "github.com/limero/offlinerss/server/newsblur/api"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Login(username, password string) error {
	args := m.Called(username, password)
	return args.Error(0)
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

func (m *MockClient) ReaderRiverStories_StoryHash(storyHash []string) ([]newsblur.ApiStory, error) {
	args := m.Called(storyHash)
	return args.Get(0).([]newsblur.ApiStory), args.Error(1)
}

func (m *MockClient) MarkStoryHashesAsRead(storyHash []string) error {
	args := m.Called(storyHash)
	return args.Error(0)
}

func (m *MockClient) MarkStoryHashAsUnread(storyHash string) error {
	args := m.Called(storyHash)
	return args.Error(0)
}

func (m *MockClient) MarkStoryHashAsStarred(storyHash string) error {
	args := m.Called(storyHash)
	return args.Error(0)
}

func (m *MockClient) MarkStoryHashAsUnstarred(storyHash string) error {
	args := m.Called(storyHash)
	return args.Error(0)
}
