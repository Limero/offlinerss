package domain

type Feed struct {
	ID      int64
	Unread  int
	Title   string
	Url     string // Link to the feed (usually .xml/.json)
	Website string // Link to the website
	Stories []*Story
}

type Feeds []*Feed

func (feeds Feeds) AddFeed(newFeed *Feed) (newFeeds Feeds) {
	/*
		Add a feed to a struct of feeds in alphabetized order
	*/

	newFeedAdded := false

	for _, feed := range feeds {
		if !newFeedAdded && newFeed.Title < feed.Title {
			newFeeds = append(newFeeds, newFeed)
			newFeedAdded = true
		}

		newFeeds = append(newFeeds, feed)
	}

	if !newFeedAdded {
		return append(newFeeds, newFeed)
	}

	return newFeeds
}

func (feeds *Feeds) GetOrCreateFeed(id int64, title string, url string, website string) *Feed {
	for _, feed := range *feeds {
		if feed.ID == id {
			return feed
		}
	}

	newFeed := &Feed{
		ID:      id,
		Unread:  0,
		Title:   title,
		Url:     url,
		Website: website,
		Stories: Stories{},
	}

	*feeds = feeds.AddFeed(newFeed)
	return newFeed
}
