package models

const (
	ActionStoryRead   = 1
	ActionStoryUnread = 2
)

type SyncToAction struct {
	Id     string
	Action int
}

type Folder struct {
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
	Timestamp string
	Hash      string
	Title     string
	Authors   string
	Content   string
	Url       string
	Unread    bool
	Date      string
}
