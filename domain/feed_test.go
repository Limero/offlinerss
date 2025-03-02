package domain

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

func TestGetOrCreateFeed(t *testing.T) {
	var feeds Feeds

	t.Run("create feed", func(t *testing.T) {
		feed := feeds.GetOrCreateFeed(1, "a", "", "")
		assert.Len(t, feeds, 1)
		assert.Equal(t, int64(1), feed.ID)
		assert.Equal(t, "a", feed.Title)
	})

	t.Run("get existing feed", func(t *testing.T) {
		feed := feeds.GetOrCreateFeed(1, "b", "", "")
		assert.Len(t, feeds, 1)
		assert.Equal(t, int64(1), feed.ID)
		assert.Equal(t, "a", feed.Title)
	})

	t.Run("create another feed", func(t *testing.T) {
		feed := feeds.GetOrCreateFeed(2, "c", "", "")
		assert.Len(t, feeds, 2)
		assert.Equal(t, int64(2), feed.ID)
		assert.Equal(t, "c", feed.Title)
	})
}
