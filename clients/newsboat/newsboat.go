package newsboat

import (
	"fmt"
	"os"
	"time"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/limero/offlinerss/helpers"
	"github.com/limero/offlinerss/models"
)

func GetChanges(clientConfig models.ClientConfig) ([]models.SyncToAction, error) {
	return helpers.GetChangesFromSqlite(
		clientConfig,
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

func GenerateCache(folders []*models.Folder, clientConfig models.ClientConfig) error {
	tmpCachePath := fmt.Sprintf("%s/cache-%d.db", os.TempDir(), time.Now().UnixNano())
	defer os.Remove(tmpCachePath)

	fmt.Println("Creating newsboat temporary cache")
	conn, err := sqlite3.Open(tmpCachePath)
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.BusyTimeout(5 * time.Second)

	fmt.Println("Creating tables in newsboat new temporary cache")
	if err := conn.Exec(`
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
			base VARCHAR(128) NOT NULL DEFAULT ""
		)`); err != nil {
		return err
	}

	newsboatUrls := ""

	fmt.Printf("Iterating over %d folders\n", len(folders))
	for _, folder := range folders {
		fmt.Printf("Iterating over %d feeds in '%s' folder\n", len(folder.Feeds), folder.Title)
		for _, feed := range folder.Feeds {
			// Newsboat stores urls in a separate file
			newsboatUrls += fmt.Sprintf("%d", feed.Id) // id instead of url to disable manual refresh
			if folder.Title != "" {
				newsboatUrls += " " + "\"" + folder.Title + "\""
			}
			newsboatUrls += "\n"

			fmt.Printf("Add feed to database: %s\n", feed.Title)
			if err := conn.Exec(
				"INSERT INTO rss_feed (rssurl, url, title) VALUES (?, ?, ?)",
				feed.Id, // id instead of url to disable manual refresh
				feed.Website,
				feed.Title,
			); err != nil {
				return err
			}

			fmt.Printf("Adding %d stories in feed %s\n", len(feed.Stories), feed.Title)
			for _, story := range feed.Stories {
				if err := conn.Exec(
					"INSERT INTO rss_item (guid, title, author, url, feedurl, pubDate, content, unread, flags) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
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

	masterCachePath, err := helpers.GetMasterCachePath(clientConfig.Type)
	if err != nil {
		return err
	}
	if err := helpers.CopyFile(tmpCachePath, masterCachePath, clientConfig.Paths.Cache); err != nil {
		return err
	}

	if err := helpers.WriteFile(newsboatUrls, clientConfig.Paths.Urls); err != nil {
		return err
	}

	return nil
}
