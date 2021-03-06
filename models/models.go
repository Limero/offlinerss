package models

const (
	ActionStoryRead      = 1
	ActionStoryUnread    = 2
	ActionStoryStarred   = 3
	ActionStoryUnstarred = 4
)

type SyncToAction struct {
	Id     string
	Action int
}

type Folder struct {
	Id    int
	Title string
	Feeds []*Feed
}

type Feed struct {
	Id      int
	Unread  int
	Title   string
	Url     string // Link to the feed (usually .xml/.json)
	Website string // Link to the website
	Stories []*Story
}

type Story struct {
	Timestamp string // Example: 1600000000
	Hash      string
	Title     string
	Authors   string
	Content   string
	Url       string
	Unread    bool
	Date      string // Example: 2006-01-02 15:04:05
	Starred   bool
}
