package client

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/limero/offlinerss/helpers"
	"github.com/limero/offlinerss/models"
)

type Newsboat struct {
	config models.ClientConfig
}

func NewNewsboat(config models.ClientConfig) *Newsboat {
	return &Newsboat{
		config: config,
	}
}

func (c Newsboat) Name() string {
	return c.config.Type
}

func (c Newsboat) GetChanges() ([]models.SyncToAction, error) {
	return helpers.GetChangesFromSqlite(
		c.config,
		"rss_item",
		"guid",
		"unread",
		"1",
		"0",
		"flags",
		"'s'",
		"''",
	)
}

func (c Newsboat) CreateNewCache() error {
	tmpCachePath := helpers.NewTmpCachePath()
	defer os.Remove(tmpCachePath)

	fmt.Println("Creating newsboat temporary cache")
	db, err := sql.Open("sqlite3", tmpCachePath)
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Println("Creating tables in newsboat new temporary cache")
	// NOT NULL constraint removed from guid so id hack in GetChangesFromSqlite works
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
			guid VARCHAR(64),
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

	masterCachePath := helpers.GetMasterCachePath(c.config.Type)

	if err := helpers.CopyFile(tmpCachePath, masterCachePath, c.config.Paths.Cache); err != nil {
		return err
	}

	return nil
}

func (c Newsboat) AddToCache(folders []*models.Folder) error {
	tmpCachePath := helpers.NewTmpCachePath()
	defer os.Remove(tmpCachePath)

	masterCachePath := helpers.GetMasterCachePath(c.config.Type)

	if err := helpers.CopyFile(masterCachePath, tmpCachePath); err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", tmpCachePath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Mark all items as read, as we might miss read events and never mark them otherwise
	if _, err = db.Exec("UPDATE rss_item SET unread = false"); err != nil {
		return err
	}

	var newsboatUrls []string

	fmt.Printf("Iterating over %d folders\n", len(folders))
	for _, folder := range folders {
		fmt.Printf("Iterating over %d feeds in '%s' folder\n", len(folder.Feeds), folder.Title)
		for _, feed := range folder.Feeds {
			// Newsboat stores urls in a separate file
			u := fmt.Sprintf("%d", feed.Id) // id instead of url to disable manual refresh
			if folder.Title != "" {
				u += " " + "\"" + folder.Title + "\""
			}
			newsboatUrls = append(newsboatUrls, u)

			fmt.Printf("Add feed to database: %s\n", feed.Title)
			if _, err = db.Exec(
				"INSERT OR REPLACE INTO rss_feed (rssurl, url, title) VALUES (?, ?, ?)",
				feed.Id, // id instead of url to disable manual refresh
				feed.Website,
				feed.Title,
			); err != nil {
				return err
			}

			fmt.Printf("Adding %d stories in feed %s\n", len(feed.Stories), feed.Title)
			for _, story := range feed.Stories {
				if _, err = db.Exec(
					"INSERT OR REPLACE INTO rss_item (guid, title, author, url, feedurl, pubDate, content, unread, flags) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
					story.Hash, // our format is different, newsboat takes the <id> in <entry> if exists
					story.Title,
					story.Authors,
					story.Url,
					feed.Id, // id instead of url to disable manual refresh
					story.Timestamp,
					story.Content,
					story.Unread,
					helpers.CondString(story.Starred, "s", ""),
				); err != nil {
					return err
				}
			}
		}
	}

	if err := helpers.CopyFile(tmpCachePath, masterCachePath, c.config.Paths.Cache); err != nil {
		return err
	}

	if err := helpers.MergeToFile(newsboatUrls, c.config.Paths.Urls); err != nil {
		return err
	}

	return nil
}