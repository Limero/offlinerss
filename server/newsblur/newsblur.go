package newsblur

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/limero/go-newsblur"
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
)

type NewsblurClient interface {
	Login(username, password string) (output *newsblur.LoginOutput, err error)

	ReaderFeeds() (output *newsblur.ReaderFeedsOutput, err error)
	ReaderUnreadStoryHashes() ([]string, error)
	ReaderStarredStoryHashes() ([]string, error)
	ReaderRiverStories_StoryHash(storyHash []string) (output *newsblur.StoriesOutput, err error)

	MarkStoryHashesAsRead(storyHash []string) (output *newsblur.MarkStoryHashesAsReadOutput, err error)
	MarkStoryHashAsUnread(storyHash string) (output *newsblur.MarkStoryHashAsUnreadOutput, err error)
	MarkStoryHashAsStarred(storyHash string) (output *newsblur.MarkStoryHashAsStarredOutput, err error)
	MarkStoryHashAsUnstarred(storyHash string) (output *newsblur.MarkStoryHashAsUnstarredOutput, err error)
}

type Newsblur struct {
	config models.ServerConfig
	client NewsblurClient
}

func New(config models.ServerConfig) *Newsblur {
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

func (s *Newsblur) Name() models.ServerName {
	return s.config.Name
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

	storyHashes, err := s.getUnreadStoryHashes()
	if err != nil {
		return nil, err
	}

	starredStoryHashes, err := s.getStarredStoryHashes()
	if err != nil {
		return nil, err
	}
	storyHashes = append(storyHashes, starredStoryHashes...)

	if err = s.fetchStories(&folders, storyHashes); err != nil {
		return nil, err
	}

	return folders, nil
}

func (s *Newsblur) getUnreadStoryHashes() ([]string, error) {
	log.Debug("Calling external NewsBlur API: ReaderUnreadStoryHashes")
	return s.client.ReaderUnreadStoryHashes()
}

func (s *Newsblur) getStarredStoryHashes() ([]string, error) {
	log.Debug("Calling external NewsBlur API: ReaderStarredStoryHashes")
	return s.client.ReaderStarredStoryHashes()
}

func (s *Newsblur) fetchStories(folders *models.Folders, storyHashes []string) error {
	var storiesOutput *newsblur.StoriesOutput
	var err error

	perPage := 100

	for page := 1; true; page++ {
		from := (page - 1) * perPage
		to := (page) * perPage
		if to > len(storyHashes) {
			to = len(storyHashes)
		}
		if from >= to {
			return nil
		}
		currentHashes := storyHashes[from:to]

		log.Debug("Calling external NewsBlur API: ReaderRiverStories. Number of storyHashes: %d. Page: %d", len(currentHashes), page)
		storiesOutput, err = s.client.ReaderRiverStories_StoryHash(currentHashes)
		if err != nil {
			return err
		}

		s.mapStoriesToFeeds(folders, storiesOutput.Stories)

		// Note that this might be fewer than the number of storyHashes
		// because ReaderRiverStories skips "disliked" intelligence trainer items
		log.Debug("Stories added: %d", len(storiesOutput.Stories))
	}
	return nil
}

func (s *Newsblur) mapStoriesToFeeds(folders *models.Folders, stories []newsblur.ApiStory) {
	for _, story := range stories {
		storyFeed := folders.FindFeed(int64(story.StoryFeedID))
		if storyFeed == nil {
			log.Debug("Could not find feed %d. Skipping story %q", story.StoryFeedID, story.StoryTitle)
			continue
		}

		// Append if latest story in feed is not the same as this one
		if len(storyFeed.Stories) == 0 || storyFeed.Stories[len(storyFeed.Stories)-1].Hash != story.StoryHash {
			unread := story.ReadStatus != 1
			if story.Starred {
				// For some reason, read starred stories show up as unread
				unread = false
			}

			storyFeed.Stories = append(storyFeed.Stories, &models.Story{
				Timestamp: time.Unix(story.StoryTimestamp, 0),
				Hash:      story.StoryHash,
				Title:     story.StoryTitle,
				Authors:   story.StoryAuthors,
				Content:   story.StoryContent,
				Url:       story.StoryPermalink,
				Unread:    unread,
				Starred:   story.Starred,
			})
		}
	}
}

func (s *Newsblur) getFolders() (models.Folders, error) {
	log.Debug("Calling external NewsBlur API: ReaderFeeds")
	readerFeedsOutput, err := s.client.ReaderFeeds()
	if err != nil {
		return nil, err
	}

	newFolders := make(models.Folders, len(readerFeedsOutput.Folders))
	for i, folder := range readerFeedsOutput.Folders {
		newFolder := models.Folder{
			Title: folder.Title,
			Feeds: models.Feeds{},
		}
		for _, feedID := range folder.FeedIDs {
			s.addFeedToFolder(readerFeedsOutput, feedID, &newFolder)
		}
		newFolders[i] = &newFolder
	}

	return newFolders, nil
}

func (s *Newsblur) addFeedToFolder(readerFeedsOutput *newsblur.ReaderFeedsOutput, feedID int, newFolder *models.Folder) {
	for _, tmpFeed := range readerFeedsOutput.Feeds {
		if feedID != tmpFeed.ID {
			continue
		}
		newFolder.Feeds = newFolder.Feeds.AddFeed(&models.Feed{
			ID:      int64(tmpFeed.ID),
			Unread:  tmpFeed.Ps + tmpFeed.Nt,
			Title:   tmpFeed.FeedTitle,
			Url:     tmpFeed.FeedAddress,
			Website: tmpFeed.FeedLink,
		})
		return
	}
}

func (s *Newsblur) SyncToServer(syncToActions models.SyncToActions) error {
	var readHashes []string
	for _, syncToAction := range syncToActions {
		switch syncToAction.Action {
		case models.ActionStoryRead:
			// Batch read events so only one request has to be done
			readHashes = append(readHashes, syncToAction.ID)
		case models.ActionStoryUnread:
			// Batching of unread events is not supported by NewsBlur, so just handle individually directly
			log.Debug("Item with hash %s has been marked as unread", syncToAction.ID)
			if err := s.markStoriesAsUnread(syncToAction.ID); err != nil {
				return err
			}
		case models.ActionStoryStarred:
			// Batching of starred events is not supported by NewsBlur, so just handle individually directly
			log.Debug("Item with hash %s has been marked as starred", syncToAction.ID)
			if err := s.markStoriesAsStarred(syncToAction.ID); err != nil {
				return err
			}
		case models.ActionStoryUnstarred:
			// Batching of unstarred events is not supported by NewsBlur, so just handle individually directly
			log.Debug("Item with hash %s has been marked as unstarred", syncToAction.ID)
			if err := s.markStoriesAsUnstarred(syncToAction.ID); err != nil {
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
