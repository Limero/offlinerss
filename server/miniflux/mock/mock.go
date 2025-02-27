package mock

import (
	"github.com/stretchr/testify/mock"
	miniflux "miniflux.app/v2/client"
)

type MockAPI struct {
	mock.Mock
}

func (m *MockAPI) Entry(entryID int64) (*miniflux.Entry, error) {
	args := m.Called(entryID)
	return args.Get(0).(*miniflux.Entry), args.Error(1)
}

func (m *MockAPI) Entries(filter *miniflux.Filter) (*miniflux.EntryResultSet, error) {
	args := m.Called(filter)
	return args.Get(0).(*miniflux.EntryResultSet), args.Error(1)
}

func (m *MockAPI) UpdateEntries(entryIDs []int64, status string) error {
	args := m.Called(entryIDs, status)
	return args.Error(0)
}

func (m *MockAPI) ToggleBookmark(entryID int64) error {
	args := m.Called(entryID)
	return args.Error(0)
}
