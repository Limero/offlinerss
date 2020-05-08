package feedreader

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
		"articles",
		"guidHash",
		"unread",
		"9",
		"8",
		"marked",
		"11",
		"10",
	)
}

func GenerateCache(folders []*models.Folder, clientConfig models.ClientConfig) error {
	tmpCachePath := fmt.Sprintf("%s/cache-%d.db", os.TempDir(), time.Now().UnixNano())
	defer os.Remove(tmpCachePath)

	fmt.Println("Creating feedreader temporary cache")
	conn, err := sqlite3.Open(tmpCachePath)
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.BusyTimeout(5 * time.Second)

	fmt.Println("Creating tables in feedreader new temporary cache")

	if err := conn.Exec(`
		CREATE TABLE "CachedActions"
			(
				"action" INTEGER NOT NULL,
				"id" TEXT NOT NULL,
				"argument" INTEGER
			);
		CREATE TABLE "Enclosures"
			(
				"articleID" TEXT NOT NULL,
				"url" TEXT NOT NULL,
				"type" INTEGER NOT NULL,
				FOREIGN KEY(articleID) REFERENCES articles(articleID)
			);
		CREATE TABLE "articles"
			(
				"articleID" TEXT PRIMARY KEY NOT NULL UNIQUE,
				"feedID" TEXT NOT NULL,
				"title" TEXT NOT NULL,
				"author" TEXT,
				"url" TEXT NOT NULL,
				"html" TEXT NOT NULL,
				"preview" TEXT NOT NULL,
				"unread" INTEGER NOT NULL,
				"marked" INTEGER NOT NULL,
				"date" INTEGER NOT NULL,
				"guidHash" TEXT,
				"lastModified" INTEGER,
				"contentFetched" INTEGER NOT NULL
			);
		CREATE TABLE "categories"
			(
				"categorieID" TEXT PRIMARY KEY NOT NULL UNIQUE,
				"title" TEXT NOT NULL,
				"orderID" INTEGER,
				"exists" INTEGER,
				"Parent" TEXT,
				"Level" INTEGER
			);
		CREATE TABLE "feeds"
			(
				"feed_id" TEXT PRIMARY KEY NOT NULL UNIQUE,
				"name" TEXT NOT NULL,
				"url" TEXT NOT NULL,
				"category_id" TEXT,
				"subscribed" INTEGER DEFAULT 1,
				"xmlURL" TEXT,
				"iconURL" TEXT
			);
		CREATE TABLE "taggings"
			(
				"articleID" TEXT NOT NULL,
				"tagID" TEXT NOT NULL,
				FOREIGN KEY(articleID) REFERENCES articles(articleID),
				FOREIGN KEY(tagID) REFERENCES tags(tagID)
			);
		CREATE TABLE "tags"
			(
				"tagID" TEXT PRIMARY KEY NOT NULL UNIQUE,
				"title" TEXT NOT NULL,
				"exists" INTEGER,
				"color" INTEGER
			)
	`); err != nil {
		return err
	}

	fmt.Printf("Iterating over %d folders\n", len(folders))
	for i, folder := range folders {
		fmt.Printf("Add folder to database: %s\n", folder.Title)
		category := 0 // 0 = Uncategorized
		if folder.Title != "" {
			category = i + 1
			if err := conn.Exec(
				"INSERT INTO categories (categorieID, title, Parent, Level) VALUES (?, ?, ?, ?)",
				category,
				folder.Title,
				-2, // ???
				1,  // ???
			); err != nil {
				return err
			}
		}

		fmt.Printf("Iterating over %d feeds in '%s' folder\n", len(folder.Feeds), folder.Title)
		for _, feed := range folder.Feeds {
			fmt.Printf("Add feed to database: %s\n", feed.Title)
			if err := conn.Exec(
				"INSERT INTO feeds (feed_id, name, url, category_id, xmlURL) VALUES (?, ?, ?, ?, ?)",
				feed.Id,
				feed.Title,
				feed.Website,
				category,
				feed.Url,
			); err != nil {
				return err
			}

			fmt.Printf("Adding %d stories in feed %s\n", len(feed.Stories), feed.Title)
			for _, story := range feed.Stories {
				if err := conn.Exec(
					"INSERT INTO articles (articleID, feedID, title, url, html, preview, unread, marked, date, guidHash, lastModified, contentFetched) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
					story.Hash,
					feed.Id,
					story.Title,
					story.Url,
					story.Content,
					story.Content,
					helpers.CondString(story.Unread, "9", "8"),
					helpers.CondString(story.Starred, "11", "10"),
					story.Timestamp,
					story.Hash,
					0,
					0,
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

	return nil
}
