package newsboat

import (
	"database/sql"
	_ "embed"
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

func New(config models.ClientConfig) *Newsboat {
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

//go:embed ddl.sql
var ddl []byte

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
	if _, err = db.Exec(string(ddl)); err != nil {
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

	if err = helpers.MarkOldStoriesAsReadAndUnstarred(db, c.GetDatabaseInfo()); err != nil {
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
