package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddFolderToFolders(t *testing.T) {
	var folders []*Folder

	feeds := []*Feed{
		&Feed{
			Id: 1,
		},
	}

	folders = AddFolderToFolders(folders, &Folder{
		Title: "cc",
		Feeds: feeds,
	})

	folders = AddFolderToFolders(folders, &Folder{
		Title: "ab",
		Feeds: feeds,
	})

	folders = AddFolderToFolders(folders, &Folder{
		Title: "bc",
		Feeds: feeds,
	})

	folders = AddFolderToFolders(folders, &Folder{
		Title: "aa",
		Feeds: feeds,
	})

	folders = AddFolderToFolders(folders, &Folder{
		Title: "ab",
		Feeds: feeds,
	})

	folders = AddFolderToFolders(folders, &Folder{
		Title: "dd",
		Feeds: feeds,
	})

	folders = AddFolderToFolders(folders, &Folder{
		Title: "",
		Feeds: feeds,
	})

	folders = AddFolderToFolders(folders, &Folder{
		Title: "nofeeds",
	})

	assert.Equal(t, 7, len(folders))
	assert.Equal(t, "", folders[0].Title)
	assert.Equal(t, "aa", folders[1].Title)
	assert.Equal(t, "ab", folders[2].Title)
	assert.Equal(t, "ab", folders[3].Title)
	assert.Equal(t, "bc", folders[4].Title)
	assert.Equal(t, "cc", folders[5].Title)
	assert.Equal(t, "dd", folders[6].Title)
}
