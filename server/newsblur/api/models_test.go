package newsblur

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReaderFeedsOutput(t *testing.T) {
	in := `
{
  "feeds": [
    {
      "id": 8312062,
      "feed_title": "OpenBSD Webzine",
      "feed_address": "https://webzine.puffy.cafe/atom.xml",
      "feed_link": "https://webzine.puffy.cafe/",
      "num_subscribers": 22,
      "updated": "1 hour",
      "updated_seconds_ago": 4161,
      "fs_size_bytes": 118793,
      "archive_count": 13,
      "last_story_date": "2023-03-19 22:00:21",
      "last_story_seconds_ago": 10177113,
      "stories_last_month": 0,
      "average_stories_per_month": 1,
      "min_to_decay": 240,
      "subs": 22,
      "is_push": false,
      "is_newsletter": false,
      "fetched_once": true,
      "search_indexed": true,
      "not_yet_fetched": false,
      "favicon_color": "c4c4c4",
      "favicon_fade": "e2e2e2",
      "favicon_border": "939393",
      "favicon_text_color": "black",
      "favicon_fetching": false,
      "favicon_url": "https://s3.amazonaws.com/icons.newsblur.com/8312062.png",
      "s3_page": false,
      "s3_icon": true,
      "disabled_page": false,
      "ps": 0,
      "nt": 0,
      "ng": 0,
      "active": true,
      "feed_opens": 4,
      "subscribed": true
    },
    {
      "id": 7600810,
      "feed_title": "OpenBSD Journal",
      "feed_address": "https://undeadly.org/cgi?action=rss",
      "feed_link": "https://undeadly.org/",
      "num_subscribers": 9,
      "updated": "2 hours",
      "updated_seconds_ago": 8244,
      "fs_size_bytes": 148598,
      "archive_count": 96,
      "last_story_date": "2023-07-14 12:19:07",
      "last_story_seconds_ago": 103187,
      "stories_last_month": 11,
      "average_stories_per_month": 5,
      "min_to_decay": 228,
      "subs": 9,
      "is_push": false,
      "is_newsletter": false,
      "fetched_once": true,
      "search_indexed": false,
      "not_yet_fetched": false,
      "favicon_color": "5c5c31",
      "favicon_fade": "7a7a4f",
      "favicon_border": "454524",
      "favicon_text_color": "white",
      "favicon_fetching": false,
      "favicon_url": "https://s3.amazonaws.com/icons.newsblur.com/7600810.png",
      "s3_page": false,
      "s3_icon": true,
      "disabled_page": false,
      "ps": 0,
      "nt": 0,
      "ng": 0,
      "active": true,
      "feed_opens": 58,
      "subscribed": true
    }
  ],
  "social_feeds": [],
  "social_profile": {},
  "social_services": {},
  "user_profile": {},
  "is_staff": false,
  "user_id": 1,
  "folders": [
    1000000,
    1111111,
    {
      "BSD": [
        7600810,
        8312062
      ]
    }
  ],
  "starred_count": 0,
  "starred_counts": [],
  "saved_searches": [],
  "dashboard_rivers": [],
  "categories": null,
  "share_ext_token": "",
  "result": "ok",
  "authenticated": true
}
`

	var raw ReaderFeedsOutputRaw

	err := json.Unmarshal([]byte(in), &raw)
	require.NoError(t, err)

	output, err := raw.toOutput()
	require.NoError(t, err)

	assert.Len(t, output.Folders, 2)
	assert.Equal(t, []Folder{
		{
			Title: "BSD",
			FeedIDs: []int{
				7600810,
				8312062,
			},
		},
		{
			Title: "",
			FeedIDs: []int{
				1000000,
				1111111,
			},
		},
	}, output.Folders)

}
