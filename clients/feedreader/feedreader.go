package feedreader

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/limero/offlinerss/helpers"
	"github.com/limero/offlinerss/models"
)

func GetChanges(clientConfig models.ClientConfig) ([]models.SyncToAction, error) {
	masterCachePath, err := helpers.GetMasterCachePath(clientConfig.Type)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(masterCachePath); os.IsNotExist(err) {
		fmt.Printf("Master cache does not exist at %s, nothing to sync to server\n", masterCachePath)
		return nil, nil
	}
	if _, err := os.Stat(clientConfig.Paths.Cache); os.IsNotExist(err) {
		fmt.Printf("Cache does not exist at %s, nothing to sync to server\n", clientConfig.Paths.Cache)
		return nil, nil
	}

	fmt.Printf("Open master cache %s\n", masterCachePath)
	conn, err := sqlite3.Open(masterCachePath)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	conn.BusyTimeout(5 * time.Second)

	// A one query workaround to get guid to also show up in sqldiff
	err = conn.Exec("UPDATE articles SET guidHash=''")
	if err != nil {
		return nil, err
	}

	diffs, err := helpers.SqlDiff(masterCachePath, clientConfig.Paths.Cache)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Iterating over %d database differences\n", len(diffs))
	var syncToActions []models.SyncToAction
	for _, row := range diffs {
		if strings.Contains(row, " unread=") {
			fmt.Printf("Change to unread: %s\n", row)
			if !strings.Contains(row, "guidHash=") {
				if err != nil {
					return nil, errors.New("guidHash not found: " + row)
				}
			}
			hash := strings.Split(strings.Split(row, "guidHash='")[1], "'")[0]
			if strings.Contains(row, " unread=8") {
				syncToActions = append(syncToActions, models.SyncToAction{
					Id:     hash,
					Action: models.ActionStoryRead,
				})
			} else if strings.Contains(row, " unread=9") {
				syncToActions = append(syncToActions, models.SyncToAction{
					Id:     hash,
					Action: models.ActionStoryUnread,
				})
			} else {
				// should never be reached
				if err != nil {
					return nil, errors.New("unread contains junk: " + row)
				}
			}
		}
	}
	return syncToActions, nil
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

	var unread int

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

			fmt.Printf("Iterating over %d stories in feed %s\n", len(feed.Stories), feed.Title)
			for _, story := range feed.Stories {
				if story.Unread {
					unread = 9
				} else {
					unread = 8
				}

				fmt.Printf("\tAdd story to database: %s\n", story.Title)
				if err := conn.Exec(
					"INSERT INTO articles (articleID, feedID, title, url, html, preview, unread, marked, date, guidHash, lastModified, contentFetched) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
					story.Hash,
					feed.Id,
					story.Title,
					story.Url,
					story.Content,
					story.Content,
					unread,
					10, // ???
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
