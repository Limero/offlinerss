package client

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/limero/offlinerss/helpers"
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
)

type Newsboat struct {
	DataPath models.DataPath
	config   models.ClientConfig
}

func NewNewsboat(config models.ClientConfig) *Newsboat {
	return &Newsboat{
		DataPath: models.GetClientDataPath(config.Type),
		config:   config,
	}
}

func (c Newsboat) Name() string {
	return "newsboat"
}

func (c Newsboat) UserDB() string {
	return c.DataPath.GetFile("cache.db")
}

func (c Newsboat) ReferenceDB() string {
	return c.DataPath.GetReferenceDB()
}

func (c Newsboat) GetChanges() (models.SyncToActions, error) {
	return helpers.GetChangesFromSqlite(
		c.ReferenceDB(),
		c.UserDB(),
		c.GetDatabaseInfo(),
	)
}

func (c Newsboat) GetDatabaseInfo() models.DatabaseInfo {
	return models.DatabaseInfo{
		StoriesTable:    "rss_item",
		StoriesIdColumn: "guid",
		Unread: models.ColumnInfo{
			Column:   "unread",
			Positive: "1",
			Negative: "0",
		},
		Starred: models.ColumnInfo{
			Column:   "flags",
			Positive: "s",
			Negative: "",
		},
	}
}

func (c Newsboat) CreateNewCache() error {
	tmpCachePath := helpers.NewTmpCachePath()
	defer os.Remove(tmpCachePath)

	log.Debug("Creating newsboat temporary cache")
	db, err := sql.Open("sqlite3", tmpCachePath)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Debug("Creating tables in newsboat new temporary cache")
	if _, err = db.Exec(`
		CREATE TABLE "rss_feed" (
			"rssurl"	VARCHAR(1024) NOT NULL,
			"url"	VARCHAR(1024) NOT NULL,
			"title"	VARCHAR(1024) NOT NULL,
			lastmodified INTEGER(11) NOT NULL DEFAULT 0,
			is_rtl INTEGER(1) NOT NULL DEFAULT 0,
			etag VARCHAR(128) NOT NULL DEFAULT "",
			PRIMARY KEY("rssurl")
		);
		CREATE TABLE rss_item (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			guid VARCHAR(64) NOT NULL,
			title VARCHAR(1024) NOT NULL,
			author VARCHAR(1024) NOT NULL,
			url VARCHAR(1024) NOT NULL,
			feedurl VARCHAR(1024) NOT NULL,
			pubDate INTEGER NOT NULL,
			content VARCHAR(65535) NOT NULL,
			unread INTEGER(1) NOT NULL ,
			enclosure_url VARCHAR(1024),
			enclosure_type VARCHAR(1024),
			enqueued INTEGER(1) NOT NULL DEFAULT 0,
			flags VARCHAR(52),
			deleted INTEGER(1) NOT NULL DEFAULT 0,
			base VARCHAR(128) NOT NULL DEFAULT "",
			UNIQUE("guid")
		)`); err != nil {
		return err
	}

	if err := helpers.CopyFile(tmpCachePath, c.ReferenceDB(), c.UserDB()); err != nil {
		return err
	}

	return nil
}

func (c Newsboat) AddToCache(folders models.Folders) error {
	tmpCachePath := helpers.NewTmpCachePath()
	defer os.Remove(tmpCachePath)

	if err := helpers.CopyFile(c.ReferenceDB(), tmpCachePath); err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", tmpCachePath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Mark all items as read, as we might miss read events and never mark them otherwise
	// Everything currently unread should be included in the folders we are adding here
	if _, err = db.Exec("UPDATE rss_item SET unread = false"); err != nil {
		return err
	}

	var newsboatUrls []string

	log.Debug("Iterating over %d folders", len(folders))
	for _, folder := range folders {
		log.Debug("Iterating over %d feeds in '%s' folder", len(folder.Feeds), folder.Title)
		for _, feed := range folder.Feeds {
			// Newsboat stores urls in a separate file
			u := fmt.Sprintf("%d", feed.Id) // id instead of url to disable manual refresh
			if folder.Title != "" {
				u += " " + "\"" + folder.Title + "\""
			}
			newsboatUrls = append(newsboatUrls, u)

			log.Debug("Add feed to database: %s", feed.Title)
			if _, err = db.Exec(
				"INSERT OR REPLACE INTO rss_feed (rssurl, url, title) VALUES (?, ?, ?)",
				feed.Id, // id instead of url to disable manual refresh
				feed.Website,
				feed.Title,
			); err != nil {
				return err
			}

			log.Debug("Adding %d stories in feed %s", len(feed.Stories), feed.Title)
			for _, story := range feed.Stories {
				if _, err = db.Exec(
					"INSERT OR REPLACE INTO rss_item (guid, title, author, url, feedurl, pubDate, content, unread, flags) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
					story.Hash, // our format is different, newsboat takes the <id> in <entry> if exists
					story.Title,
					story.Authors,
					story.Url,
					feed.Id, // id instead of url to disable manual refresh
					story.Timestamp.Unix(),
					story.Content,
					story.Unread,
					helpers.CondString(story.Starred, "s", ""),
				); err != nil {
					return err
				}
			}
		}
	}

	if err := helpers.CopyFile(tmpCachePath, c.ReferenceDB(), c.UserDB()); err != nil {
		return err
	}

	if err := helpers.MergeToFile(
		newsboatUrls,
		c.DataPath.GetFile("urls"),
		urlsSortFunc(),
	); err != nil {
		return err
	}

	return nil
}

func urlsSortFunc() func(s1, s2 string) bool {
	return func(s1, s2 string) bool {
		// Split lines into words
		words1 := strings.Fields(s1)
		words2 := strings.Fields(s2)

		// Compare the last words
		return words1[len(words1)-1] < words2[len(words2)-1]
	}
}
