package models

type Feed struct {
	Id      int
	Unread  int
	Title   string
	Url     string // Link to the feed (usually .xml/.json)
	Website string // Link to the website
	Stories []*Story
}

type Feeds []*Feed

func AddFeedToFeeds(feeds []*Feed, newFeed *Feed) (newFeeds []*Feed) {
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
