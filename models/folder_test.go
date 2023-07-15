package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddFolder(t *testing.T) {
	var folders Folders

	feeds := Feeds{
		{
			Id: 1,
		},
	}

	folders = folders.AddFolder(&Folder{
		Title: "cc",
		Feeds: feeds,
	})

	folders = folders.AddFolder(&Folder{
		Title: "aa",
		Feeds: feeds,
	})

	folders = folders.AddFolder(&Folder{
		Title: "bb",
		Feeds: feeds,
	})

	folders = folders.AddFolder(&Folder{
		Title: "",
		Feeds: feeds,
	})

	folders = folders.AddFolder(&Folder{
		Title: "nofeeds",
	})

	assert.Equal(t, 5, len(folders))
	assert.Equal(t, "", folders[0].Title)
	assert.Equal(t, "aa", folders[1].Title)
	assert.Equal(t, "bb", folders[2].Title)
	assert.Equal(t, "cc", folders[3].Title)
	assert.Equal(t, "nofeeds", folders[4].Title)
}

func TestGetOrCreateFolder(t *testing.T) {
	var folders Folders

	t.Run("create folder", func(t *testing.T) {
		folder := folders.GetOrCreateFolder(1, "a")
		assert.Len(t, folders, 1)
		assert.Equal(t, int64(1), folder.Id)
		assert.Equal(t, "a", folder.Title)
	})

	t.Run("get existing folder", func(t *testing.T) {
		folder := folders.GetOrCreateFolder(1, "b")
		assert.Len(t, folders, 1)
		assert.Equal(t, int64(1), folder.Id)
		assert.Equal(t, "a", folder.Title)
	})

	t.Run("create another folder", func(t *testing.T) {
		folder := folders.GetOrCreateFolder(2, "c")
		assert.Len(t, folders, 2)
		assert.Equal(t, int64(2), folder.Id)
		assert.Equal(t, "c", folder.Title)
	})
}
