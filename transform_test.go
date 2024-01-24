package main

import (
	"testing"

	"github.com/limero/offlinerss/models"
	"github.com/stretchr/testify/assert"
)

func TestTransformFolders(t *testing.T) {
	folders := models.Folders{
		{
			Feeds: models.Feeds{
				{
					Stories: models.Stories{
						{
							Content: "a https://example.com/?utm_source=rss b",
							Url:     "https://example.com/?utm_source=rss",
						},
					},
				},
			},
		},
	}

	TransformFolders(folders)
	/*
		assert.Equal(t,
			"a https://example.com/ b",
			folders[0].Feeds[0].Stories[0].Content,
		)
	*/
	assert.Equal(t,
		"https://example.com/",
		folders[0].Feeds[0].Stories[0].Url,
	)
}
