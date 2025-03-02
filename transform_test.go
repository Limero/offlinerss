package main

import (
	"testing"

	"github.com/limero/offlinerss/domain"
	"github.com/stretchr/testify/assert"
)

func TestTransformFolders(t *testing.T) {
	folders := domain.Folders{
		{
			Feeds: domain.Feeds{
				{
					Stories: domain.Stories{
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
	assert.Equal(t,
		"a https://example.com/ b",
		folders[0].Feeds[0].Stories[0].Content,
	)
	assert.Equal(t,
		"https://example.com/",
		folders[0].Feeds[0].Stories[0].Url,
	)
}
