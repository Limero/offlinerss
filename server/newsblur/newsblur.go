package newsblur

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/limero/offlinerss/domain"
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/server/newsblur/api"
)

type API interface {
	Login(username, password string) error

	ReaderFeeds() (output *api.ReaderFeedsOutput, err error)
	ReaderUnreadStoryHashes() ([]string, error)
	ReaderStarredStoryHashes() ([]api.HashWithTimestamp, error)
	ReaderRiverStories_StoryHash(storyHash []string) ([]api.Story, error)

	MarkStoryHashesAsRead(storyHash []string) error
	MarkStoryHashAsUnread(storyHash []string) error
	MarkStoryHashAsStarred(storyHash []string) error
	MarkStoryHashAsUnstarred(storyHash []string) error
}

type Newsblur struct {
	config domain.ServerConfig
	api    API
}

func New(config domain.ServerConfig) *Newsblur {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	client := api.New(&http.Client{
		Jar: cookieJar,
	})
	if config.Hostname != "" {
		client.Hostname = config.Hostname
	}
	return &Newsblur{
		config: config,
		api:    client,
	}
}

func (s *Newsblur) Name() domain.ServerName {
	return s.config.Name
}

func (s *Newsblur) Login() error {
	log.Debug("Calling external NewsBlur API: Login")
	return s.api.Login(s.config.Username, s.config.Password)
}

func (s *Newsblur) GetFoldersWithStories(from *time.Time) (domain.Folders, error) {
	folders, err := s.getFolders()
	if err != nil {
		return nil, err
	}

	log.Debug("Calling external NewsBlur API: ReaderUnreadStoryHashes")
	storyHashes, err := s.api.ReaderUnreadStoryHashes()
	if err != nil {
		return nil, err
	}

	log.Debug("Calling external NewsBlur API: ReaderStarredStoryHashes")
	starredStoryHashes, err := s.api.ReaderStarredStoryHashes()
	if err != nil {
		return nil, err
	}
	for _, hash := range starredStoryHashes {
		if from != nil && from.After(time.Unix(hash.Timestamp, 0).UTC()) {
			// The newest hashes are in the beginning, so we can stop once we find one older
			break
		}
		storyHashes = append(storyHashes, hash.Hash)
	}

	if err = s.fetchStories(&folders, storyHashes); err != nil {
		return nil, err
	}

	return folders, nil
}

func (s *Newsblur) fetchStories(folders *domain.Folders, storyHashes []string) error {
	var stories []api.Story
	var err error

	perPage := 100

	for page := 1; true; page++ {
		from := (page - 1) * perPage
		to := min((page)*perPage, len(storyHashes))
		if from >= to {
			return nil
		}
		currentHashes := storyHashes[from:to]

		log.Debug("Calling external NewsBlur API: ReaderRiverStories. Number of storyHashes: %d. Page: %d", len(currentHashes), page)
		stories, err = s.api.ReaderRiverStories_StoryHash(currentHashes)
		if err != nil {
			return err
		}

		s.mapStoriesToFeeds(folders, stories)

		// Note that this might be fewer than the number of storyHashes
		// because ReaderRiverStories skips "disliked" intelligence trainer items
		log.Debug("Stories added: %d", len(stories))
	}
	return nil
}

func (s *Newsblur) mapStoriesToFeeds(folders *domain.Folders, stories []api.Story) {
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

			storyFeed.Stories = append(storyFeed.Stories, &domain.Story{
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

func (s *Newsblur) getFolders() (domain.Folders, error) {
	log.Debug("Calling external NewsBlur API: ReaderFeeds")
	readerFeedsOutput, err := s.api.ReaderFeeds()
	if err != nil {
		return nil, err
	}

	newFolders := make(domain.Folders, len(readerFeedsOutput.Folders))
	for i, folder := range readerFeedsOutput.Folders {
		newFolder := domain.Folder{
			Title: folder.Title,
			Feeds: domain.Feeds{},
		}
		for _, feedID := range folder.FeedIDs {
			s.addFeedToFolder(readerFeedsOutput, feedID, &newFolder)
		}
		newFolders[i] = &newFolder
	}

	return newFolders, nil
}

func (s *Newsblur) addFeedToFolder(readerFeedsOutput *api.ReaderFeedsOutput, feedID int, newFolder *domain.Folder) {
	for _, tmpFeed := range readerFeedsOutput.Feeds {
		if feedID != tmpFeed.ID {
			continue
		}
		newFolder.Feeds = newFolder.Feeds.AddFeed(&domain.Feed{
			ID:      int64(tmpFeed.ID),
			Unread:  tmpFeed.Ps + tmpFeed.Nt,
			Title:   tmpFeed.FeedTitle,
			Url:     tmpFeed.FeedAddress,
			Website: tmpFeed.FeedLink,
		})
		return
	}
}

func (s *Newsblur) MarkStoriesAsRead(hashes []string) error {
	log.Debug("Calling external NewsBlur API: MarkStoryHashesAsRead. Hashes: %+v", hashes)
	return s.api.MarkStoryHashesAsRead(hashes)
}

func (s *Newsblur) MarkStoriesAsUnread(hashes []string) error {
	log.Debug("Calling external NewsBlur API: MarkStoryHashAsUnread. Hashes: %+v", hashes)
	return s.api.MarkStoryHashAsUnread(hashes)
}

func (s *Newsblur) MarkStoriesAsStarred(hashes []string) error {
	log.Debug("Calling external NewsBlur API: MarkStoryHashAsStarred. Hash: %+v", hashes)
	return s.api.MarkStoryHashAsStarred(hashes)
}

func (s *Newsblur) MarkStoriesAsUnstarred(hashes []string) error {
	log.Debug("Calling external NewsBlur API: MarkStoryHashAsUnstarred. Hashes: %+v", hashes)
	return s.api.MarkStoryHashAsUnstarred(hashes)
}
