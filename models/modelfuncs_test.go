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

func Test_AddFeedToFeeds(t *testing.T) {
	var feeds Feeds

	feeds = AddFeedToFeeds(feeds, &Feed{
		Title: "cc",
	})
	feeds = AddFeedToFeeds(feeds, &Feed{
		Title: "aa",
	})
	feeds = AddFeedToFeeds(feeds, &Feed{
		Title: "bb",
	})

	assert.Equal(t, 3, len(feeds))
	assert.Equal(t, "aa", feeds[0].Title)
	assert.Equal(t, "bb", feeds[1].Title)
	assert.Equal(t, "cc", feeds[2].Title)
}
