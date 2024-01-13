package main

import "github.com/limero/offlinerss/models"

func TransformFolders(folders models.Folders) {
	for _, folder := range folders {
		for _, feed := range folder.Feeds {
			for _, story := range feed.Stories {
				story.Url = transformURL(story.Url)
			}
		}
	}
}

func transformURL(u string) string {
	// TODO
	return u
}
