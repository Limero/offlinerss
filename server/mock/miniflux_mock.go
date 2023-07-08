package mock

import (
	"github.com/stretchr/testify/mock"
	miniflux "miniflux.app/client"
)

type MockMinifluxClient struct {
	mock.Mock
}

func (m *MockMinifluxClient) Entry(entryID int64) (*miniflux.Entry, error) {
	args := m.Called(entryID)
	return args.Get(0).(*miniflux.Entry), args.Error(1)
}

func (m *MockMinifluxClient) Entries(filter *miniflux.Filter) (*miniflux.EntryResultSet, error) {
	args := m.Called(filter)
	return args.Get(0).(*miniflux.EntryResultSet), args.Error(1)
}

func (m *MockMinifluxClient) UpdateEntries(entryIDs []int64, status string) error {
	args := m.Called(entryIDs, status)
	return args.Error(0)
}

func (m *MockMinifluxClient) ToggleBookmark(entryID int64) error {
	args := m.Called(entryID)
	return args.Error(0)
}
