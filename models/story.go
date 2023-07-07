package models

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

type Stories []*Story
