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

func Test_AddFeedToFeeds(t *testing.T) {
	var feeds []*Feed

	feeds = AddFeedToFeeds(feeds, &Feed{
		Title: "cc",
	})
	feeds = AddFeedToFeeds(feeds, &Feed{
		Title: "ab",
	})
	feeds = AddFeedToFeeds(feeds, &Feed{
		Title: "bc",
	})
	feeds = AddFeedToFeeds(feeds, &Feed{
		Title: "aa",
	})
	feeds = AddFeedToFeeds(feeds, &Feed{
		Title: "ab",
	})
	feeds = AddFeedToFeeds(feeds, &Feed{
		Title: "dd",
	})

	assert.Equal(t, 6, len(feeds))
	assert.Equal(t, "aa", feeds[0].Title)
	assert.Equal(t, "ab", feeds[1].Title)
	assert.Equal(t, "ab", feeds[2].Title)
	assert.Equal(t, "bc", feeds[3].Title)
	assert.Equal(t, "cc", feeds[4].Title)
	assert.Equal(t, "dd", feeds[5].Title)
}
