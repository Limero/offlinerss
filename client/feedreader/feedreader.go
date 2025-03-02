package feedreader

import (
	_ "embed"

	"github.com/limero/offlinerss/client"
	"github.com/limero/offlinerss/domain"
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/util"
)

type Feedreader struct {
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

func New(config domain.ClientConfig) *Feedreader {
	return &Feedreader{
		client.Client{
			ClientName: domain.ClientFeedReader,
			DataPath:   domain.GetClientDataPath(config.Name),
			Config:     config,
			DatabaseInfo: domain.DatabaseInfo{
				FileName:        "feedreader-7.db",
				DDL:             ddl,
				StoriesTable:    "articles",
				StoriesIDColumn: "guidHash",
				Unread: domain.ColumnInfo{
					Column:   "unread",
					Positive: "9",
					Negative: "8",
				},
				Starred: domain.ColumnInfo{
					Column:   "marked",
					Positive: "11",
					Negative: "10",
				},
			},
			Files: domain.ClientFiles{
				{
					FileName: "feedreader-7.db",
					TargetPaths: []string{
						util.DataDir("feedreader/data/feedreader-7.db"),
					},
				},
			},
		},
	}
}

func (c Feedreader) AddToCache(folders domain.Folders) error {
	tmpCachePath, db, closer, err := c.CreateNewTmpCache()
	defer closer()
	if err != nil {
		return err
	}

	log.Debug("Iterating over %d folders", len(folders))
	for i, folder := range folders {
		log.Debug("Add folder to database: %s", folder.Title)
		category := 0 // 0 = Uncategorized
		if folder.Title != "" {
			category = i + 1
			if _, err = db.Exec(
				// TODO: categorieID is unique and if folders are changed between syncs, this might mess up
				insertFolder,
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
				insertFeed,
				feed.ID,
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
					insertStory,
					story.Hash,
					feed.ID,
					story.Title,
					story.Url,
					story.Content,
					story.Content,
					util.Cond(story.Unread, "9", "8"),
					util.Cond(story.Starred, "11", "10"),
					story.Timestamp.Unix(),
					story.Hash,
					0,
					0,
				); err != nil {
					return err
				}
			}
		}
	}

	return util.CopyFile(tmpCachePath, c.ReferenceDB(), c.UserDB())
}
