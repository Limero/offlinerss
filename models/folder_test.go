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
