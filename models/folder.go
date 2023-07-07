package models

type Folder struct {
	Id    int
	Title string
	Feeds []*Feed
}

type Folders []*Folder
