package mock

import (
	newsblur "github.com/limero/offlinerss/server/newsblur/api"
	"github.com/stretchr/testify/mock"
)

type MockAPI struct {
	mock.Mock
}

func (m *MockAPI) Login(username, password string) error {
	args := m.Called(username, password)
	return args.Error(0)
}

func (m *MockAPI) ReaderFeeds() (*newsblur.ReaderFeedsOutput, error) {
	args := m.Called()
	return args.Get(0).(*newsblur.ReaderFeedsOutput), args.Error(1)
}

func (m *MockAPI) ReaderUnreadStoryHashes() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAPI) ReaderStarredStoryHashes() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAPI) ReaderRiverStories_StoryHash(storyHash []string) ([]newsblur.Story, error) {
	args := m.Called(storyHash)
	return args.Get(0).([]newsblur.Story), args.Error(1)
}

func (m *MockAPI) MarkStoryHashesAsRead(storyHash []string) error {
	args := m.Called(storyHash)
	return args.Error(0)
}

func (m *MockAPI) MarkStoryHashAsUnread(storyHash string) error {
	args := m.Called(storyHash)
	return args.Error(0)
}

func (m *MockAPI) MarkStoryHashAsStarred(storyHash string) error {
	args := m.Called(storyHash)
	return args.Error(0)
}

func (m *MockAPI) MarkStoryHashAsUnstarred(storyHash string) error {
	args := m.Called(storyHash)
	return args.Error(0)
}
