package apimodels

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReaderStarredStoryHashes(t *testing.T) {
	in := `
{
  "starred_story_hashes": [
    "1111111:aaaaaa",
    "2222222:bbbbbb"
  ],
  "result": "ok",
  "authenticated": true,
  "user_id": 500000
}
`

	var output ReaderStarredStoryHashes

	err := json.Unmarshal([]byte(in), &output)
	require.NoError(t, err)

	assert.Len(t, output.StarredStoryHashes, 2)
}

func TestReaderUnreadStoryHashes(t *testing.T) {
	in := `
{
  "unread_feed_story_hashes": {
    "1111111": [
      "1111111:aaaaaa"
    ],
    "2222222": [
      "2222222:bbbbbb",
      "3333333:cccccc"
    ]
  },
  "result": "ok",
  "authenticated": true,
  "user_id": 500000
}`

	var api ReaderUnreadStoryHashes

	err := json.Unmarshal([]byte(in), &api)
	require.NoError(t, err)

	output := api.ToOutput()

	assert.Len(t, output, 3)
	assert.Equal(t, []string{
		"1111111:aaaaaa",
		"2222222:bbbbbb",
		"3333333:cccccc",
	}, output)
}
