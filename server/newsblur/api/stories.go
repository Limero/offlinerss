package newsblur

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/limero/offlinerss/server/newsblur/api/apimodels"
)

// Retrieve stories from a single feed.
// GET /reader/feed/:id
// https://www.newsblur.com/api#/reader/feed/:id
func (nb *Newsblur) ReaderFeed(feedID string, page int) (output *ReaderFeedOutput, err error) {
	if page == 0 {
		page = 1
	}
	body, err := GetWithBody(
		nb.client,
		fmt.Sprintf("%s/reader/feed/%s?page=%d", nb.Hostname, feedID, page),
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// Retrieve a user's starred stories.
// GET /reader/starred_stories
// https://newsblur.com/api#/reader/starred_stories
func (nb *Newsblur) ReaderStarredStories(page int) (output *StoriesOutput, err error) {
	if page == 0 {
		page = 1
	}

	body, err := GetWithBody(
		nb.client,
		fmt.Sprintf("%s/reader/starred_stories?page=%d", nb.Hostname, page),
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// Retrieve the story hashes of a user's starred stories.
// GET /reader/starred_story_hashes
// https://newsblur.com/api#/reader/starred_story_hashes
func (nb *Newsblur) ReaderStarredStoryHashes() ([]string, error) {
	body, err := GetWithBody(
		nb.client,
		nb.Hostname+"/reader/starred_story_hashes",
	)
	if err != nil {
		return nil, err
	}

	var output *apimodels.ReaderStarredStoryHashes
	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}

	return output.StarredStoryHashes, nil
}

// Retrieve stories from a collection of feeds
// GET /reader/river_stories
// https://www.newsblur.com/api#/reader/river_stories
func (nb *Newsblur) ReaderRiverStories(feeds []string, page int) (output *StoriesOutput, err error) {
	if page == 0 {
		page = 1
	}

	formData := url.Values{
		"feeds": feeds,
		"page":  {strconv.Itoa(page)},
	}

	body, err := PostWithBody(nb.client, nb.Hostname+"/reader/river_stories", formData)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// Retrieve up to 100 stories when specifying by story_hash.
// GET /reader/river_stories
// https://newsblur.com/api#/reader/river_stories
func (nb *Newsblur) ReaderRiverStories_StoryHash(storyHash []string) (output *StoriesOutput, err error) {
	formData := url.Values{
		"h": storyHash,
	}

	body, err := PostWithBody(nb.client, nb.Hostname+"/reader/river_stories", formData)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// The story_hashes of all unread stories.
// GET /reader/unread_story_hashes
// https://newsblur.com/api#/reader/unread_story_hashes
func (nb *Newsblur) ReaderUnreadStoryHashes() ([]string, error) {
	body, err := GetWithBody(
		nb.client,
		nb.Hostname+"/reader/unread_story_hashes",
	)
	if err != nil {
		return nil, err
	}

	var output *apimodels.ReaderUnreadStoryHashes
	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}

	return output.ToOutput(), nil
}

// Mark stories as read using their unique story_hash.
// POST /reader/mark_story_hashes_as_read
// https://www.newsblur.com/api#/reader/mark_story_hashes_as_read
func (nb *Newsblur) MarkStoryHashesAsRead(storyHash []string) (output *MarkStoryHashesAsReadOutput, err error) {
	formData := url.Values{
		"story_hash": storyHash,
	}

	body, err := PostWithBody(nb.client, nb.Hostname+"/reader/mark_story_hashes_as_read", formData)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// Mark a single story as unread using its unique story_hash.
// POST /reader/mark_story_hash_as_unread
// https://www.newsblur.com/api#/reader/mark_story_hash_as_unread
func (nb *Newsblur) MarkStoryHashAsUnread(storyHash string) (output *MarkStoryHashAsUnreadOutput, err error) {
	formData := url.Values{
		"story_hash": {storyHash},
	}

	body, err := PostWithBody(nb.client, nb.Hostname+"/reader/mark_story_hash_as_unread", formData)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// Mark a story as starred (saved).
// POST /reader/mark_story_hash_as_starred
// https://www.newsblur.com/api#/reader/mark_story_hash_as_starred
func (nb *Newsblur) MarkStoryHashAsStarred(storyHash string) (output *MarkStoryHashAsStarredOutput, err error) {
	formData := url.Values{
		"story_hash": {storyHash},
	}

	body, err := PostWithBody(nb.client, nb.Hostname+"/reader/mark_story_hash_as_starred", formData)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// Mark a story as unstarred (unsaved).
// POST /reader/mark_story_hash_as_unstarred
// https://www.newsblur.com/api#/reader/mark_story_hash_as_unstarred
func (nb *Newsblur) MarkStoryHashAsUnstarred(storyHash string) (output *MarkStoryHashAsUnstarredOutput, err error) {
	formData := url.Values{
		"story_hash": {storyHash},
	}

	body, err := PostWithBody(nb.client, nb.Hostname+"/reader/mark_story_hash_as_unstarred", formData)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}

	return output, nil
}
