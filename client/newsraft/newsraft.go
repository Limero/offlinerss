package newsraft

import (
	_ "embed"
	"strings"

	"github.com/limero/offlinerss/client"
	"github.com/limero/offlinerss/domain"
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/util"
)

type Newsraft struct {
	client.Client
}

//go:embed sql/ddl.sql
var ddl string

//go:embed sql/insert-feed.sql
var insertFeed string

//go:embed sql/insert-story.sql
var insertStory string

func New(config domain.ClientConfig) *Newsraft {
	return &Newsraft{
		client.Client{
			ClientName: domain.ClientNewsraft,
			DataPath:   domain.GetClientDataPath(config.Name),
			Config:     config,
			DatabaseInfo: domain.DatabaseInfo{
				FileName:        "newsraft.sqlite3",
				DDL:             ddl,
				StoriesTable:    "items",
				StoriesIDColumn: "guid",
				Unread: domain.ColumnInfo{
					Column:   "unread",
					Positive: "1",
					Negative: "0",
				},
				Starred: domain.ColumnInfo{
					Column:   "important",
					Positive: "1",
					Negative: "0",
				},
			},
			Files: domain.ClientFiles{
				{
					FileName: "newsraft.sqlite3",
					TargetPaths: []string{
						util.DataDir("newsraft/newsraft.sqlite3"),
					},
				},
				{
					FileName: "feeds",
					TargetPaths: []string{
						util.ConfigDir("newsraft/feeds"),
					},
				},
			},
		},
	}
}

const fs = "\x1F"

func (c Newsraft) AddToCache(folders domain.Folders) error {
	tmpCachePath, db, closer, err := c.CreateNewTmpCache()
	defer closer()
	if err != nil {
		return err
	}

	var newsraftUrls []string

	log.Debug("Iterating over %d folders", len(folders))
	for _, folder := range folders {
		log.Debug("Iterating over %d feeds in '%s' folder", len(folder.Feeds), folder.Title)
		for _, feed := range folder.Feeds {
			if !strings.HasPrefix(feed.Url, "http") {
				log.Debug("Skip url %s because it doesn't start with http", feed.Url)
				continue
			}

			// Newsraft stores urls in a separate file
			u := feed.Url + " \""
			if folder.Title != "" {
				u += "[" + folder.Title + "] "
			}
			u += feed.Title + "\""
			newsraftUrls = append(newsraftUrls, u)

			log.Debug("Add feed to database: %s", feed.Title)
			if _, err = db.Exec(
				insertFeed,
				feed.Url,
				feed.Title,
				feed.Website,
			); err != nil {
				return err
			}

			log.Debug("Adding %d stories in feed %s", len(feed.Stories), feed.Title)
			for _, story := range feed.Stories {
				if _, err = db.Exec(
					insertStory,
					feed.Url,
					story.Hash,
					story.Title,
					story.Url,
					fs+"type=text/html"+fs+"text="+story.Content, // TODO: dynamic type
					fs+"type=author"+fs+"name="+story.Authors,
					story.Timestamp.Unix(),
					story.Unread,
					story.Starred,
				); err != nil {
					return err
				}
			}
		}
	}

	if err := util.MergeToFile(
		newsraftUrls,
		c.DataPath.GetFile("feeds"),
		nil,
	); err != nil {
		return err
	}

	return util.CopyFile(tmpCachePath, c.ReferenceDB(), c.UserDB())
}
