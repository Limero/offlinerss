package client

import (
	"database/sql"
	"os"

	"github.com/limero/offlinerss/helpers"
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
)

type Feedreader struct {
	DataPath models.DataPath
	config   models.ClientConfig
}

func NewFeedreader(config models.ClientConfig) *Feedreader {
	return &Feedreader{
		DataPath: models.GetClientDataPath(config.Type),
		config:   config,
	}
}

func (c Feedreader) Name() string {
	return "feedreader"
}

func (c Feedreader) UserDB() string {
	return c.DataPath.GetFile("feedreader-7.db")
}

func (c Feedreader) ReferenceDB() string {
	return c.DataPath.GetReferenceDB()
}

func (c Feedreader) GetChanges() (models.SyncToActions, error) {
	return helpers.GetChangesFromSqlite(
		c.ReferenceDB(),
		c.UserDB(),
		c.GetDatabaseInfo(),
	)
}

func (c Feedreader) GetDatabaseInfo() models.DatabaseInfo {
	return models.DatabaseInfo{
		StoriesTable:    "articles",
		StoriesIdColumn: "guidHash",
		Unread: models.ColumnInfo{
			Column:   "unread",
			Positive: "9",
			Negative: "8",
		},
		Starred: models.ColumnInfo{
			Column:   "marked",
			Positive: "11",
			Negative: "10",
		},
	}
}

func (c Feedreader) CreateNewCache() error {
	tmpCachePath := helpers.NewTmpCachePath()
	defer os.Remove(tmpCachePath)

	log.Debug("Creating feedreader temporary cache")
	db, err := sql.Open("sqlite3", tmpCachePath)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Debug("Creating tables in feedreader new temporary cache")

	if _, err = db.Exec(`
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

	if err := helpers.CopyFile(tmpCachePath, c.ReferenceDB(), c.UserDB()); err != nil {
		return err
	}

	return nil
}

func (c Feedreader) AddToCache(folders models.Folders) error {
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

	log.Debug("Iterating over %d folders", len(folders))
	for i, folder := range folders {
		log.Debug("Add folder to database: %s", folder.Title)
		category := 0 // 0 = Uncategorized
		if folder.Title != "" {
			category = i + 1
			if _, err = db.Exec(
				// TODO: categorieID is unique and if folders are changed between syncs, this might mess up
				"INSERT OR REPLACE INTO categories (categorieID, title, Parent, Level) VALUES (?, ?, ?, ?)",
				category,
				folder.Title,
				-2, // ???
				1,  // ???
			); err != nil {
				return err
			}
		}

		log.Debug("Iterating over %d feeds in '%s' folder", len(folder.Feeds), folder.Title)
		for _, feed := range folder.Feeds {
			log.Debug("Add feed to database: %s", feed.Title)
			if _, err = db.Exec(
				"INSERT OR REPLACE INTO feeds (feed_id, name, url, category_id, xmlURL) VALUES (?, ?, ?, ?, ?)",
				feed.Id,
				feed.Title,
				feed.Website,
				category,
				feed.Url,
			); err != nil {
				return err
			}

			log.Debug("Adding %d stories in feed %s", len(feed.Stories), feed.Title)
			for _, story := range feed.Stories {
				if _, err = db.Exec(
					"INSERT OR REPLACE INTO articles (articleID, feedID, title, url, html, preview, unread, marked, date, guidHash, lastModified, contentFetched) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
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

	if err := helpers.CopyFile(tmpCachePath, c.ReferenceDB(), c.UserDB()); err != nil {
		return err
	}

	return nil
}
