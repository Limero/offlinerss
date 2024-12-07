package apimodels

type ReaderStarredStoryHashes struct {
	StarredStoryHashes []string `json:"starred_story_hashes"`
	Result             string   `json:"result"`
	Authenticated      bool     `json:"authenticated"`
	UserID             int      `json:"user_id"`
}

type ReaderUnreadStoryHashes struct {
	UnreadFeedStoryHashes map[string][]string `json:"unread_feed_story_hashes"`
	Result                string              `json:"result"`
	Authenticated         bool                `json:"authenticated"`
	UserID                int                 `json:"user_id"`
}

func (api ReaderUnreadStoryHashes) ToOutput() []string {
	var storyHashes []string

	for _, hashes := range api.UnreadFeedStoryHashes {
		storyHashes = append(storyHashes, hashes...)
	}

	return storyHashes
}
