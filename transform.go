package main

import (
	"github.com/limero/linkcleaner"
	"github.com/limero/offlinerss/domain"
	"github.com/limero/offlinerss/log"
)

func TransformFolders(folders domain.Folders) {
	for _, folder := range folders {
		for _, feed := range folder.Feeds {
			for _, story := range feed.Stories {
				story.Content = transformContent(story.Content)
				story.Url = transformURL(story.Url)
			}
		}
	}
}

func transformContent(c string) string {
	return linkcleaner.CleanAllURLsInString(c)
}

func transformURL(u string) string {
	cleanURL, err := linkcleaner.CleanURLString(u)
	if err != nil {
		log.Warn("Failed to clean url %q: %v", u, err)
		return u
	}
	return cleanURL.String()
}
