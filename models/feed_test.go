package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
