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
							Url: "https://example.com/?utm_source=rss",
						},
					},
				},
			},
		},
	}

	TransformFolders(folders)
	assert.Equal(t,
		"https://example.com/",
		folders[0].Feeds[0].Stories[0].Url,
	)
}
