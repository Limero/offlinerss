package newsblur

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Newsblur struct {
	Hostname string
	client   *http.Client
}

func New(client *http.Client) *Newsblur {
	return &Newsblur{
		Hostname: "https://www.newsblur.com",
		client:   client,
	}
}

// Login as an existing user.
// POST /api/login
// https://www.newsblur.com/api#/api/login
func (nb *Newsblur) Login(username, password string) error {
	resp, err := nb.client.PostForm(nb.Hostname+"/api/login", url.Values{
		"username": {username},
		"password": {password},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var output struct {
		Authenticated bool `json:"authenticated"`
		Errors        any  `json:"errors"`
	}

	if err := json.Unmarshal(body, &output); err != nil {
		return err
	}

	if !output.Authenticated {
		return fmt.Errorf("failed to login to newsblur: %v", output.Errors)
	}

	return nil
}

// Retrieve a list of feeds to which a user is actively subscribed.
// GET /reader/feeds
// https://www.newsblur.com/api#/reader/feeds
func (nb *Newsblur) ReaderFeeds() (output *ReaderFeedsOutput, err error) {
	resp, err := nb.client.Get(nb.Hostname + "/reader/feeds?v=2")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw struct {
		Folders []any     `json:"folders"`
		Feeds   []ApiFeed `json:"feeds"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	output = &ReaderFeedsOutput{
		Folders: make([]Folder, 0),
		Feeds:   raw.Feeds,
	}

	emptyFolder := Folder{
		Title:   "",
		FeedIDs: []int{},
	}

	for _, element := range raw.Folders {
		switch element.(type) {
		case float64, float32:
			// Feed without folder
			emptyFolder.FeedIDs = append(emptyFolder.FeedIDs, int(element.(float64)))
		case map[string]any:
			// Feed with folder
			folders := element.(map[string]any)
			for folder, feeds := range folders {
				feedIDs := []int{}
				for _, feedId := range feeds.([]any) {
					feedIDs = append(feedIDs, int(feedId.(float64)))
				}

				// Add folder if it's not empty
				if len(feedIDs) > 0 {
					output.Folders = append(output.Folders, Folder{
						Title:   folder,
						FeedIDs: feedIDs,
					})
				}
			}
		}
	}

	if len(emptyFolder.FeedIDs) > 0 {
		output.Folders = append(output.Folders, emptyFolder)
	}

	return output, nil
}

// Retrieve a user's starred stories.
// GET /reader/starred_stories
// https://newsblur.com/api#/reader/starred_stories
func (nb *Newsblur) ReaderStarredStories(page int) ([]ApiStory, error) {
	if page == 0 {
		page = 1
	}

	resp, err := nb.client.Get(fmt.Sprintf("%s/reader/starred_stories?page=%d", nb.Hostname, page))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var output struct {
		Stories []ApiStory `json:"stories"`
	}

	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}

	return output.Stories, nil
}

// Retrieve the story hashes of a user's starred stories.
// GET /reader/starred_story_hashes
// https://newsblur.com/api#/reader/starred_story_hashes
func (nb *Newsblur) ReaderStarredStoryHashes() ([]string, error) {
	resp, err := nb.client.Get(nb.Hostname + "/reader/starred_story_hashes")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var output struct {
		StarredStoryHashes []string `json:"starred_story_hashes"`
	}
	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}

	return output.StarredStoryHashes, nil
}

// Retrieve stories from a collection of feeds
// GET /reader/river_stories
// https://www.newsblur.com/api#/reader/river_stories
func (nb *Newsblur) ReaderRiverStories(feeds []string, page int) ([]ApiStory, error) {
	if page == 0 {
		page = 1
	}

	formData := url.Values{
		"feeds": feeds,
		"page":  {strconv.Itoa(page)},
	}

	resp, err := nb.client.PostForm(nb.Hostname+"/reader/river_stories", formData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var output struct {
		Stories []ApiStory `json:"stories"`
	}

	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}

	return output.Stories, nil
}

// Retrieve up to 100 stories when specifying by story_hash.
// GET /reader/river_stories
// https://newsblur.com/api#/reader/river_stories
func (nb *Newsblur) ReaderRiverStories_StoryHash(storyHash []string) ([]ApiStory, error) {
	formData := url.Values{
		"h": storyHash,
	}

	resp, err := nb.client.PostForm(nb.Hostname+"/reader/river_stories", formData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var output struct {
		Stories []ApiStory `json:"stories"`
	}

	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}

	return output.Stories, nil
}

// The story_hashes of all unread stories.
// GET /reader/unread_story_hashes
// https://newsblur.com/api#/reader/unread_story_hashes
func (nb *Newsblur) ReaderUnreadStoryHashes() ([]string, error) {
	resp, err := nb.client.Get(nb.Hostname + "/reader/unread_story_hashes")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var output struct {
		UnreadFeedStoryHashes map[string][]string `json:"unread_feed_story_hashes"`
	}
	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}

	var storyHashes []string
	for _, hashes := range output.UnreadFeedStoryHashes {
		storyHashes = append(storyHashes, hashes...)
	}

	return storyHashes, nil
}

// Mark stories as read using their unique story_hash.
// POST /reader/mark_story_hashes_as_read
// https://www.newsblur.com/api#/reader/mark_story_hashes_as_read
func (nb *Newsblur) MarkStoryHashesAsRead(storyHash []string) error {
	_, err := nb.client.PostForm(nb.Hostname+"/reader/mark_story_hashes_as_read", url.Values{
		"story_hash": storyHash,
	})
	return err
}

// Mark a single story as unread using its unique story_hash.
// POST /reader/mark_story_hash_as_unread
// https://www.newsblur.com/api#/reader/mark_story_hash_as_unread
func (nb *Newsblur) MarkStoryHashAsUnread(storyHash string) error {
	_, err := nb.client.PostForm(nb.Hostname+"/reader/mark_story_hash_as_unread", url.Values{
		"story_hash": {storyHash},
	})
	return err
}

// Mark a story as starred (saved).
// POST /reader/mark_story_hash_as_starred
// https://www.newsblur.com/api#/reader/mark_story_hash_as_starred
func (nb *Newsblur) MarkStoryHashAsStarred(storyHash string) error {
	_, err := nb.client.PostForm(nb.Hostname+"/reader/mark_story_hash_as_starred", url.Values{
		"story_hash": {storyHash},
	})
	return err
}

// Mark a story as unstarred (unsaved).
// POST /reader/mark_story_hash_as_unstarred
// https://www.newsblur.com/api#/reader/mark_story_hash_as_unstarred
func (nb *Newsblur) MarkStoryHashAsUnstarred(storyHash string) error {
	_, err := nb.client.PostForm(nb.Hostname+"/reader/mark_story_hash_as_unstarred", url.Values{
		"story_hash": {storyHash},
	})
	return err
}
