package quiterss

import (
	"database/sql"
	_ "embed"
	"os"

	"github.com/limero/offlinerss/helpers"
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
)

type QuiteRSS struct {
	DataPath models.DataPath
	config   models.ClientConfig
}

func New(config models.ClientConfig) *QuiteRSS {
	return &QuiteRSS{
		DataPath: models.GetClientDataPath(config.Type),
		config:   config,
	}
}

func (c QuiteRSS) Name() string {
	return "quiterss"
}

func (c QuiteRSS) UserDB() string {
	return c.DataPath.GetFile("feeds.db")
}

func (c QuiteRSS) ReferenceDB() string {
	return c.DataPath.GetReferenceDB()
}

func (c QuiteRSS) GetChanges() (models.SyncToActions, error) {
	return helpers.GetChangesFromSqlite(
		c.ReferenceDB(),
		c.UserDB(),
		c.GetDatabaseInfo(),
	)
}

func (c QuiteRSS) GetDatabaseInfo() models.DatabaseInfo {
	return models.DatabaseInfo{
		StoriesTable:    "news",
		StoriesIdColumn: "guid",
		Unread: models.ColumnInfo{
			Column:   "read",
			Positive: "0",
			Negative: "2",
		},
		Starred: models.ColumnInfo{
			Column:   "starred",
			Positive: "1",
			Negative: "0",
		},
	}
}

//go:embed ddl.sql
var ddl []byte

func (c QuiteRSS) CreateNewCache() error {
	tmpCachePath := helpers.NewTmpCachePath()
	defer os.Remove(tmpCachePath)

	log.Debug("Creating QuiteRSS temporary cache")
	db, err := sql.Open("sqlite3", tmpCachePath)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Debug("Creating tables in QuiteRSS new temporary cache")

	if _, err = db.Exec(string(ddl)); err != nil {
		return err
	}

	if err := helpers.CopyFile(tmpCachePath, c.ReferenceDB(), c.UserDB()); err != nil {
		return err
	}

	return nil
}

func (c QuiteRSS) AddToCache(folders models.Folders) error {
	// TODO: Remove this once db has been confirmed idempotent, like the Newsboat client
	if err := c.CreateNewCache(); err != nil {
		return err
	}

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

	if err = helpers.MarkOldStoriesAsReadAndUnstarred(db, c.GetDatabaseInfo()); err != nil {
		return err
	}

	latestFeedId := 0 // This is required because folder/feed share same table and use ids

	log.Debug("Iterating over %d folders", len(folders))
	for _, folder := range folders {
		log.Debug("Add folder to database: %s", folder.Title)
		category := 0 // Category variable separate to lastFeedId to support feeds without a folder
		if folder.Title != "" {
			latestFeedId++
			category = latestFeedId
			if _, err = db.Exec(
				"INSERT INTO feeds (id, text) VALUES (?, ?)",
				latestFeedId,
				folder.Title,
			); err != nil {
				return err
			}
		}

		log.Debug("Iterating over %d feeds in '%s' folder", len(folder.Feeds), folder.Title)
		for _, feed := range folder.Feeds {
			log.Debug("Add feed to database: %s", feed.Title)
			latestFeedId++
			if _, err = db.Exec(
				"INSERT INTO feeds (id, text, title, xmlUrl, htmlUrl, unread, parentId) VALUES (?, ?, ?, ?, ?, ?, ?)",
				latestFeedId,
				feed.Title,
				feed.Title,
				feed.Url,
				feed.Website,
				feed.Unread,
				category,
			); err != nil {
				return err
			}

			log.Debug("Adding %d stories in feed %s", len(feed.Stories), feed.Title)
			for _, story := range feed.Stories {
				if _, err = db.Exec(
					"INSERT INTO news (feedId, guid, description, title, published, read, starred, link_href) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
					latestFeedId,
					story.Hash,
					story.Content,
					story.Title,
					story.Timestamp.Unix(),
					helpers.CondString(story.Unread, "0", "2"),
					story.Starred,
					story.Url,
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
