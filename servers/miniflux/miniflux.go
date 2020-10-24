package miniflux

import (
	"fmt"
	"strconv"

	"github.com/limero/offlinerss/models"
	miniflux "miniflux.app/client"
)

func Login(username string, password string) (*miniflux.Client, error) {
	client := miniflux.New("http://localhost", username, password)

	if _, err := client.Me(); err != nil {
		return nil, err
	}

	return client, nil
}

func GetFoldersWithStories(client *miniflux.Client) ([]*models.Folder, error) {
	var folders []*models.Folder

	entries, err := client.Entries(&miniflux.Filter{
		Status: "unread",
	})
	if err != nil {
		return nil, err
	}

	for _, entry := range entries.Entries {
		unread := true
		if entry.Status == "read" {
			unread = false
		}

		story := &models.Story{
			Timestamp: strconv.FormatInt(entry.Date.Unix(), 10),
			Hash:      entry.Hash,
			Title:     entry.Title,
			Authors:   entry.Author,
			Content:   entry.Content,
			Url:       entry.URL,
			Unread:    unread,
			Date:      entry.Date.Format("2006-01-02 15:04:05"),
			Starred:   entry.Starred,
		}

		var storyFolder *models.Folder

		for _, folder := range folders {
			if int64(folder.Id) == entry.Feed.Category.ID {
				storyFolder = folder
				break
			}
		}
		if storyFolder == nil {
			// New folder
			storyFolder = &models.Folder{
				Id:    int(entry.Feed.Category.ID),
				Title: entry.Feed.Category.Title,
				Feeds: []*models.Feed{},
			}

			folders = models.AddFolderToFolders(folders, storyFolder)
		}

		var storyFeed *models.Feed
		for _, feed := range storyFolder.Feeds {
			if int64(feed.Id) == entry.Feed.ID {
				storyFeed = feed
				break
			}
		}
		if storyFeed == nil {
			// New feed
			storyFeed = &models.Feed{
				Id:      int(entry.Feed.ID),
				Unread:  0,
				Title:   entry.Feed.Title,
				Url:     entry.Feed.FeedURL,
				Website: entry.Feed.SiteURL,
				Stories: []*models.Story{},
			}
			storyFolder.Feeds = models.AddFeedToFeeds(storyFolder.Feeds, storyFeed)
		}

		if story.Unread {
			storyFeed.Unread++
		}

		storyFeed.Stories = append(storyFeed.Stories, story)
	}

	return folders, nil
}

func SyncToServer(client *miniflux.Client, syncToActions []models.SyncToAction) error {
	fmt.Println("----- Sync back to Miniflux is not yet supported.")
	return nil
}
