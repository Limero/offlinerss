package newsboat

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/limero/offlinerss/client"
	"github.com/limero/offlinerss/helpers"
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
)

type Newsboat struct {
	client.Client
}

//go:embed ddl.sql
var ddl []byte

func New(config models.ClientConfig) *Newsboat {
	return &Newsboat{
		client.Client{
			ClientName: "newsboat",
			DataPath:   models.GetClientDataPath(config.Type),
			Config:     config,
			DatabaseInfo: models.DatabaseInfo{
				FileName:        "cache.db",
				DDL:             ddl,
				StoriesTable:    "rss_item",
				StoriesIDColumn: "guid",
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
			},
		},
	}
}

func (c Newsboat) AddToCache(folders models.Folders) error {
	tmpCachePath, db, closer, err := c.CreateNewTmpCache()
	defer closer()
	if err != nil {
		return err
	}

	var newsboatUrls []string

	log.Debug("Iterating over %d folders", len(folders))
	for _, folder := range folders {
		log.Debug("Iterating over %d feeds in '%s' folder", len(folder.Feeds), folder.Title)
		for _, feed := range folder.Feeds {
			// Newsboat stores urls in a separate file
			u := fmt.Sprintf("%d", feed.ID) // id instead of url to disable manual refresh
			if folder.Title != "" {
				u += " " + "\"" + folder.Title + "\""
			}
			newsboatUrls = append(newsboatUrls, u)

			log.Debug("Add feed to database: %s", feed.Title)
			if _, err = db.Exec(
				"INSERT OR REPLACE INTO rss_feed (rssurl, url, title) VALUES (?, ?, ?)",
				feed.ID, // id instead of url to disable manual refresh
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
					feed.ID, // id instead of url to disable manual refresh
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

	if err := helpers.MergeToFile(
		newsboatUrls,
		c.DataPath.GetFile("urls"),
		urlsSortFunc(),
	); err != nil {
		return err
	}

	return helpers.CopyFile(tmpCachePath, c.ReferenceDB(), c.UserDB())
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
