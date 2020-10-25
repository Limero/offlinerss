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
		Status: miniflux.EntryStatusUnread,
	})
	if err != nil {
		return nil, err
	}

	for _, entry := range entries.Entries {
		unread := true
		if entry.Status == miniflux.EntryStatusRead {
			unread = false
		}

		story := &models.Story{
			Timestamp: strconv.FormatInt(entry.Date.Unix(), 10),
			Hash:      strconv.FormatInt(entry.ID, 10), // Miniflux has "hash" but IDs are used for marking entries
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

func handleStarred(client *miniflux.Client, syncToAction models.SyncToAction) error {
	// Because Miniflux only support toggling starred instead of setting it directly,
	// we have to check its current status

	actionId, err := strconv.ParseInt(syncToAction.Id, 10, 64)
	if err != nil {
		return err
	}

	entry, err := client.Entry(actionId)
	if err != nil {
		return err
	}

	if (entry.Starred && syncToAction.Action == models.ActionStoryUnstarred) ||
		(!entry.Starred && syncToAction.Action == models.ActionStoryStarred) {
		return client.ToggleBookmark(actionId)
	}

	return nil
}

func SyncToServer(client *miniflux.Client, syncToActions []models.SyncToAction) error {
	var readIds []int64
	var unreadIds []int64

	for _, syncToAction := range syncToActions {
		actionId, err := strconv.ParseInt(syncToAction.Id, 10, 64)
		if err != nil {
			return err
		}

		switch syncToAction.Action {
		case models.ActionStoryRead:
			// Batch read events so only one request has to be done
			readIds = append(readIds, actionId)
		case models.ActionStoryUnread:
			// Batch unread events so only one request has to be done
			unreadIds = append(unreadIds, actionId)
		case models.ActionStoryStarred:
			if err := handleStarred(client, syncToAction); err != nil {
				return err
			}
		case models.ActionStoryUnstarred:
			if err := handleStarred(client, syncToAction); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unsupported Miniflux syncToAction: %d", syncToAction.Action)
		}
	}

	if len(readIds) > 0 {
		if err := client.UpdateEntries(readIds, miniflux.EntryStatusRead); err != nil {
			return err
		}
		fmt.Printf("%d items has been marked as read\n", len(readIds))
	}

	if len(unreadIds) > 0 {
		if err := client.UpdateEntries(unreadIds, miniflux.EntryStatusUnread); err != nil {
			return err
		}
		fmt.Printf("%d items has been marked as unread\n", len(unreadIds))
	}

	return nil
}
