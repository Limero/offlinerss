package newsblur

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strconv"

	"github.com/limero/go-newsblur"
	"github.com/limero/offlinerss/models"
)

func Login(username string, password string) (*http.Client, error) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: cookieJar,
	}

	fmt.Println("Calling external NewsBlur API: Login")
	loginOutput, err := newsblur.ApiLogin(client, &newsblur.LoginInput{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	if !loginOutput.Authenticated {
		return nil, errors.New(fmt.Sprintf("Failed to login to NewsBlur. %v", loginOutput.Errors))
	}

	return client, nil
}

func GetFolders(client *http.Client) ([]*models.Folder, error) {
	fmt.Println("Calling external NewsBlur API: ReaderFeeds")
	readerFeedsOutput, err := newsblur.ApiReaderFeeds(client)
	if err != nil {
		return nil, err
	}

	// noFolder is a collection of feeds without folder
	noFolder := models.Folder{
		Title: "",
		Feeds: []*models.Feed{},
	}

	var newFolders []*models.Folder
	for _, element := range readerFeedsOutput.Folders {
		switch element.(type) {
		case float64, float32:
			// Feed without folder
			addFeedToFolder(readerFeedsOutput, element, &noFolder)
		case map[string]interface{}:
			// Feed with folder
			folders := element.(map[string]interface{})
			for folder, feeds := range folders {
				newFolder := models.Folder{
					Title: folder,
					Feeds: []*models.Feed{},
				}

				for _, feedId := range feeds.([]interface{}) {
					addFeedToFolder(readerFeedsOutput, feedId, &newFolder)
				}

				// Add folder if it's not empty
				if len(newFolder.Feeds) > 0 {
					newFolders = models.AddFolderToFolders(newFolders, &newFolder)
				}
			}
		}
	}

	return models.AddFolderToFolders(newFolders, &noFolder), nil
}

func addFeedToFolder(readerFeedsOutput *newsblur.ReaderFeedsOutput, feedId interface{}, newFolder *models.Folder) {
	// Loop through list of feeds to find one with matching id
	for _, tmpFeed := range readerFeedsOutput.Feeds {
		if int(feedId.(float64)) == tmpFeed.ID {
			// Match found
			if tmpFeed.Ps != 0 || tmpFeed.Nt != 0 {
				// Feed has unread items, add it
				newFolder.Feeds = models.AddFeedToFeeds(newFolder.Feeds, &models.Feed{
					Id:      tmpFeed.ID,
					Unread:  tmpFeed.Ps + tmpFeed.Nt,
					Title:   tmpFeed.FeedTitle,
					Url:     tmpFeed.FeedAddress,
					Website: tmpFeed.FeedLink,
				})
			}
			return
		}
	}
}

func GetFoldersWithStories(client *http.Client) ([]*models.Folder, error) {
	// Like GetFolders but it will also load all unread stories with it
	folders, err := GetFolders(client)
	if err != nil {
		return nil, err
	}

	var feedIds []string

	for _, folder := range folders {
		for _, feed := range folder.Feeds {
			feedIds = append(feedIds, strconv.Itoa(feed.Id))
		}
	}

	for page := 1; true; page++ {
		fmt.Printf("Calling external NewsBlur API: ReaderRiverStories. Number of feeds: %d. Page: %d\n", len(feedIds), page)
		readerRiverStoriesOutput, err := newsblur.ApiReaderRiverStories(client, &newsblur.ReaderRiverStoriesInput{
			Feeds: feedIds,
			Page:  strconv.Itoa(page),
		})
		if err != nil {
			return nil, err
		}

		// Map stories to feeds
	STORIES:
		for _, story := range readerRiverStoriesOutput.Stories {
			for _, folder := range folders {
				for _, feed := range folder.Feeds {
					if feed.Id == story.StoryFeedID {
						// Append if latest story in feed is not the same as this one
						if len(feed.Stories) == 0 || feed.Stories[len(feed.Stories)-1].Hash != story.StoryHash {
							feed.Stories = append(feed.Stories, &models.Story{
								Timestamp: story.StoryTimestamp,
								Hash:      story.StoryHash,
								Title:     story.StoryTitle,
								Authors:   story.StoryAuthors,
								Content:   story.StoryContent,
								Url:       story.StoryPermalink,
								Unread:    story.ReadStatus != 1,
								Date:      story.StoryDate,
								Starred:   story.Starred,
							})
						}
						continue STORIES
					}
				}
			}
		}

		fmt.Printf("Stories added: %d\n", len(readerRiverStoriesOutput.Stories))
		if len(readerRiverStoriesOutput.Stories) == 0 {
			break
		}
	}

	return folders, nil
}

func MarkStoriesAsRead(client *http.Client, hashes ...string) error {
	if len(hashes) == 0 {
		return nil
	}

	fmt.Printf("Calling external NewsBlur API: MarkStoryHashesAsRead. Hashes: %+v\n", hashes)
	_, err := newsblur.ApiMarkStoryHashesAsRead(client, &newsblur.MarkStoryHashesAsReadInput{
		StoryHash: hashes,
	})
	return err
}

func MarkStoriesAsUnread(client *http.Client, hashes ...string) error {
	// NewsBlur doesn't support batching unread events. So we have to handle them individually
	for _, hash := range hashes {
		fmt.Printf("Calling external NewsBlur API: MarkStoryHashAsUnread. Hash: %s\n", hash)
		_, err := newsblur.ApiMarkStoryHashAsUnread(client, &newsblur.MarkStoryHashAsUnreadInput{
			StoryHash: hash,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func MarkStoriesAsStarred(client *http.Client, hashes ...string) error {
	// NewsBlur doesn't support batching starred events. So we have to handle them individually
	for _, hash := range hashes {
		fmt.Printf("Calling external NewsBlur API: MarkStoryHashAsStarred. Hash: %s\n", hash)
		_, err := newsblur.ApiMarkStoryHashAsStarred(client, &newsblur.MarkStoryHashAsStarredInput{
			StoryHash: hash,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func MarkStoriesAsUnstarred(client *http.Client, hashes ...string) error {
	// NewsBlur doesn't support batching unstarred events. So we have to handle them individually
	for _, hash := range hashes {
		fmt.Printf("Calling external NewsBlur API: MarkStoryHashAsUnstarred. Hash: %s\n", hash)
		_, err := newsblur.ApiMarkStoryHashAsUnstarred(client, &newsblur.MarkStoryHashAsUnstarredInput{
			StoryHash: hash,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func SyncToServer(client *http.Client, syncToActions []models.SyncToAction) error {
	var readHashes []string
	for _, syncToAction := range syncToActions {
		switch syncToAction.Action {
		case models.ActionStoryRead:
			// Batch read events so only one request has to be done
			readHashes = append(readHashes, syncToAction.Id)
		case models.ActionStoryUnread:
			// Batching of unread events is not supported by NewsBlur, so just handle individually directly
			fmt.Printf("Item with hash %s has been marked as unread\n", syncToAction.Id)
			if err := MarkStoriesAsUnread(client, syncToAction.Id); err != nil {
				return err
			}
		case models.ActionStoryStarred:
			// Batching of starred events is not supported by NewsBlur, so just handle individually directly
			fmt.Printf("Item with hash %s has been marked as starred\n", syncToAction.Id)
			if err := MarkStoriesAsStarred(client, syncToAction.Id); err != nil {
				return err
			}
		case models.ActionStoryUnstarred:
			// Batching of unstarred events is not supported by NewsBlur, so just handle individually directly
			fmt.Printf("Item with hash %s has been marked as unstarred\n", syncToAction.Id)
			if err := MarkStoriesAsUnstarred(client, syncToAction.Id); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unsupported Newsblur syncToAction: %d", syncToAction.Action)
		}
	}

	if err := MarkStoriesAsRead(client, readHashes...); err != nil {
		return err
	}
	fmt.Printf("%d items has been marked as read\n", len(readHashes))
	return nil
}
