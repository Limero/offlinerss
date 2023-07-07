package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddFeed(t *testing.T) {
	var feeds Feeds

	feeds = feeds.AddFeed(&Feed{
		Title: "cc",
	})
	feeds = feeds.AddFeed(&Feed{
		Title: "aa",
	})
	feeds = feeds.AddFeed(&Feed{
		Title: "bb",
	})

	assert.Equal(t, 3, len(feeds))
	assert.Equal(t, "aa", feeds[0].Title)
	assert.Equal(t, "bb", feeds[1].Title)
	assert.Equal(t, "cc", feeds[2].Title)
}
