package server

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"time"

	"github.com/limero/go-newsblur"
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
)

type NewsblurClient interface {
	Login(username, password string) (output *newsblur.LoginOutput, err error)
	ReaderRiverStories(feeds []string, page int) (output *newsblur.ReaderRiverStoriesOutput, err error)
	ReaderFeeds() (output *newsblur.ReaderFeedsOutput, err error)
	MarkStoryHashesAsRead(storyHash []string) (output *newsblur.MarkStoryHashesAsReadOutput, err error)
	MarkStoryHashAsUnread(storyHash string) (output *newsblur.MarkStoryHashAsUnreadOutput, err error)
	MarkStoryHashAsStarred(storyHash string) (output *newsblur.MarkStoryHashAsStarredOutput, err error)
	MarkStoryHashAsUnstarred(storyHash string) (output *newsblur.MarkStoryHashAsUnstarredOutput, err error)
}

type Newsblur struct {
	config models.ServerConfig
	client NewsblurClient
}

func NewNewsblur(config models.ServerConfig) *Newsblur {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	return &Newsblur{
		config: config,
		client: newsblur.New(&http.Client{
			Jar: cookieJar,
		}),
	}
}

func (s *Newsblur) Name() string {
	return s.config.Type
}

func (s *Newsblur) Login() error {
	log.Debug("Calling external NewsBlur API: Login")
	_, err := s.client.Login(s.config.Username, s.config.Password)
	return err
}

func (s *Newsblur) GetFoldersWithStories() (models.Folders, error) {
	folders, err := s.getFolders()
	if err != nil {
		return nil, err
	}

	var feedIds []string

	for _, folder := range folders {
		for _, feed := range folder.Feeds {
			feedIds = append(feedIds, strconv.FormatInt(feed.Id, 10))
		}
	}

	for page := 1; true; page++ {
		log.Debug("Calling external NewsBlur API: ReaderRiverStories. Number of feeds: %d. Page: %d", len(feedIds), page)
		readerRiverStoriesOutput, err := s.client.ReaderRiverStories(feedIds, page)
		if err != nil {
			return nil, err
		}

		// Map stories to feeds
		for _, story := range readerRiverStoriesOutput.Stories {
			var storyFeed *models.Feed
			for _, folder := range folders {
				for _, feed := range folder.Feeds {
					if feed.Id == int64(story.StoryFeedID) {
						storyFeed = feed
					}
				}
			}

			// Append if latest story in feed is not the same as this one
			if len(storyFeed.Stories) == 0 || storyFeed.Stories[len(storyFeed.Stories)-1].Hash != story.StoryHash {
				storyFeed.Stories = append(storyFeed.Stories, &models.Story{
					Timestamp: time.Unix(story.StoryTimestamp, 0),
					Hash:      story.StoryHash,
					Title:     story.StoryTitle,
					Authors:   story.StoryAuthors,
					Content:   story.StoryContent,
					Url:       story.StoryPermalink,
					Unread:    story.ReadStatus != 1,
					Starred:   story.Starred,
				})
			}
		}

		log.Debug("Stories added: %d", len(readerRiverStoriesOutput.Stories))
		if len(readerRiverStoriesOutput.Stories) == 0 {
			break
		}
	}

	return folders, nil
}

func (s *Newsblur) getFolders() (models.Folders, error) {
	log.Debug("Calling external NewsBlur API: ReaderFeeds")
	readerFeedsOutput, err := s.client.ReaderFeeds()
	if err != nil {
		return nil, err
	}

	var newFolders models.Folders
	for _, folder := range readerFeedsOutput.Folders {
		newFolder := models.Folder{
			Title: folder.Title,
			Feeds: models.Feeds{},
		}
		for _, feedId := range folder.FeedIDs {
			s.addFeedToFolder(readerFeedsOutput, feedId, &newFolder)
		}
		newFolders = append(newFolders, &newFolder)
	}

	return newFolders, nil
}

func (s *Newsblur) addFeedToFolder(readerFeedsOutput *newsblur.ReaderFeedsOutput, feedId int, newFolder *models.Folder) {
	// Loop through list of feeds to find one with matching id
	for _, tmpFeed := range readerFeedsOutput.Feeds {
		if feedId == tmpFeed.ID {
			// Match found
			if tmpFeed.Ps != 0 || tmpFeed.Nt != 0 {
				// Feed has unread items, add it
				newFolder.Feeds = newFolder.Feeds.AddFeed(&models.Feed{
					Id:      int64(tmpFeed.ID),
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

func (s *Newsblur) SyncToServer(syncToActions models.SyncToActions) error {
	var readHashes []string
	for _, syncToAction := range syncToActions {
		switch syncToAction.Action {
		case models.ActionStoryRead:
			// Batch read events so only one request has to be done
			readHashes = append(readHashes, syncToAction.Id)
		case models.ActionStoryUnread:
			// Batching of unread events is not supported by NewsBlur, so just handle individually directly
			log.Debug("Item with hash %s has been marked as unread", syncToAction.Id)
			if err := s.markStoriesAsUnread(syncToAction.Id); err != nil {
				return err
			}
		case models.ActionStoryStarred:
			// Batching of starred events is not supported by NewsBlur, so just handle individually directly
			log.Debug("Item with hash %s has been marked as starred", syncToAction.Id)
			if err := s.markStoriesAsStarred(syncToAction.Id); err != nil {
				return err
			}
		case models.ActionStoryUnstarred:
			// Batching of unstarred events is not supported by NewsBlur, so just handle individually directly
			log.Debug("Item with hash %s has been marked as unstarred", syncToAction.Id)
			if err := s.markStoriesAsUnstarred(syncToAction.Id); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unsupported Newsblur syncToAction: %d", syncToAction.Action)
		}
	}

	if err := s.markStoriesAsRead(readHashes...); err != nil {
		return err
	}
	log.Debug("%d items has been marked as read", len(readHashes))
	return nil
}

func (s *Newsblur) markStoriesAsRead(hashes ...string) error {
	if len(hashes) == 0 {
		return nil
	}

	log.Debug("Calling external NewsBlur API: MarkStoryHashesAsRead. Hashes: %+v", hashes)
	_, err := s.client.MarkStoryHashesAsRead(hashes)
	return err
}

func (s *Newsblur) markStoriesAsUnread(hashes ...string) error {
	// NewsBlur doesn't support batching unread events. So we have to handle them individually
	for _, hash := range hashes {
		log.Debug("Calling external NewsBlur API: MarkStoryHashAsUnread. Hash: %s", hash)
		_, err := s.client.MarkStoryHashAsUnread(hash)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Newsblur) markStoriesAsStarred(hashes ...string) error {
	// NewsBlur doesn't support batching starred events. So we have to handle them individually
	for _, hash := range hashes {
		log.Debug("Calling external NewsBlur API: MarkStoryHashAsStarred. Hash: %s", hash)
		_, err := s.client.MarkStoryHashAsStarred(hash)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Newsblur) markStoriesAsUnstarred(hashes ...string) error {
	// NewsBlur doesn't support batching unstarred events. So we have to handle them individually
	for _, hash := range hashes {
		log.Debug("Calling external NewsBlur API: MarkStoryHashAsUnstarred. Hash: %s", hash)
		_, err := s.client.MarkStoryHashAsUnstarred(hash)
		if err != nil {
			return err
		}
	}
	return nil
}
