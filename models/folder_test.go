package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddFolderToFolders(t *testing.T) {
	var folders Folders

	feeds := Feeds{
		{
			Id: 1,
		},
	}

	folders = AddFolderToFolders(folders, &Folder{
		Title: "cc",
		Feeds: feeds,
	})

	folders = AddFolderToFolders(folders, &Folder{
		Title: "aa",
		Feeds: feeds,
	})

	folders = AddFolderToFolders(folders, &Folder{
		Title: "bb",
		Feeds: feeds,
	})

	folders = AddFolderToFolders(folders, &Folder{
		Title: "",
		Feeds: feeds,
	})

	folders = AddFolderToFolders(folders, &Folder{
		Title: "nofeeds",
	})

	assert.Equal(t, 5, len(folders))
	assert.Equal(t, "", folders[0].Title)
	assert.Equal(t, "aa", folders[1].Title)
	assert.Equal(t, "bb", folders[2].Title)
	assert.Equal(t, "cc", folders[3].Title)
	assert.Equal(t, "nofeeds", folders[4].Title)
}
