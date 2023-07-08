package models

type ColumnInfo struct {
	Column   string
	Positive string
	Negative string
}

type DatabaseInfo struct {
	StoriesTable    string
	StoriesIdColumn string
	Unread          ColumnInfo
	Starred         ColumnInfo
}
