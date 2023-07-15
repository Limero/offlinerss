package models

import "time"

type Story struct {
	Timestamp time.Time
	Hash      string
	Title     string
	Authors   string
	Content   string
	Url       string
	Unread    bool
	Starred   bool
}

type Stories []*Story
