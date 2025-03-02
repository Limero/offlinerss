package models

type ColumnInfo struct {
	Column   string
	Positive string
	Negative string
}

type DatabaseInfo struct {
	FileName        string
	DDL             string
	StoriesTable    string
	StoriesIDColumn string
	Unread          ColumnInfo
	Starred         ColumnInfo
}
