package quiterss

import (
	_ "embed"

	"github.com/limero/offlinerss/client"
	"github.com/limero/offlinerss/domain"
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/util"
)

type QuiteRSS struct {
	client.Client
}

//go:embed sql/ddl.sql
var ddl string

//go:embed sql/insert-folder.sql
var insertFolder string

//go:embed sql/insert-feed.sql
var insertFeed string

//go:embed sql/insert-story.sql
var insertStory string

func New(config domain.ClientConfig) *QuiteRSS {
	return &QuiteRSS{
		client.Client{
			ClientName: domain.ClientQuiteRSS,
			DataPath:   domain.GetClientDataPath(config.Name),
			Config:     config,
			DatabaseInfo: domain.DatabaseInfo{
				FileName:        "feeds.db",
				DDL:             ddl,
				StoriesTable:    "news",
				StoriesIDColumn: "guid",
				Unread: domain.ColumnInfo{
					Column:   "read",
					Positive: "0",
					Negative: "2",
				},
				Starred: domain.ColumnInfo{
					Column:   "starred",
					Positive: "1",
					Negative: "0",
				},
			},
			Files: domain.ClientFiles{
				{
					FileName: "feeds.db",
					TargetPaths: []string{
						util.DataDir("QuiteRss/QuiteRss/feeds.db"),
					},
				},
			},
		},
	}
}

func (c QuiteRSS) AddToCache(folders domain.Folders) error {
	// TODO: Remove this once db has been confirmed idempotent, like the Newsboat client
	if err := c.CreateNewCache(); err != nil {
		return err
	}

	tmpCachePath, db, closer, err := c.CreateNewTmpCache()
	defer closer()
	if err != nil {
		return err
	}

	latestFeedID := 0 // This is required because folder/feed share same table and use ids

	log.Debug("Iterating over %d folders", len(folders))
	for _, folder := range folders {
		log.Debug("Add folder to database: %s", folder.Title)
		category := 0 // Category variable separate to lastFeedId to support feeds without a folder
		if folder.Title != "" {
			latestFeedID++
			category = latestFeedID
			if _, err = db.Exec(
				insertFolder,
				latestFeedID,
				folder.Title,
			); err != nil {
				return err
			}
		}

		log.Debug("Iterating over %d feeds in '%s' folder", len(folder.Feeds), folder.Title)
		for _, feed := range folder.Feeds {
			log.Debug("Add feed to database: %s", feed.Title)
			latestFeedID++
			if _, err = db.Exec(
				insertFeed,
				latestFeedID,
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
					insertStory,
					latestFeedID,
					story.Hash,
					story.Content,
					story.Title,
					story.Timestamp.Unix(),
					util.Cond(story.Unread, "0", "2"),
					story.Starred,
					story.Url,
				); err != nil {
					return err
				}
			}
		}
	}

	return util.CopyFile(tmpCachePath, c.ReferenceDB(), c.UserDB())
}
