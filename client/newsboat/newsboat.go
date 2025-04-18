package newsboat

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/limero/offlinerss/client"
	"github.com/limero/offlinerss/domain"
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/util"
)

type Newsboat struct {
	client.Client
}

//go:embed sql/ddl.sql
var ddl string

//go:embed sql/insert-feed.sql
var insertFeed string

//go:embed sql/insert-story.sql
var insertStory string

func New(config domain.ClientConfig) *Newsboat {
	return &Newsboat{
		client.Client{
			ClientName: domain.ClientNewsboat,
			DataPath:   domain.GetClientDataPath(config.Name),
			Config:     config,
			DatabaseInfo: domain.DatabaseInfo{
				FileName:        "cache.db",
				DDL:             ddl,
				StoriesTable:    "rss_item",
				StoriesIDColumn: "guid",
				Unread: domain.ColumnInfo{
					Column:   "unread",
					Positive: "1",
					Negative: "0",
				},
				Starred: domain.ColumnInfo{
					Column:   "flags",
					Positive: "s",
					Negative: "",
				},
			},
			Files: domain.ClientFiles{
				{
					FileName: "cache.db",
					TargetPaths: []string{
						util.DataDir("newsboat/cache.db"),
					},
				},
				{
					FileName: "urls",
					TargetPaths: []string{
						util.ConfigDir("newsboat/urls"),
					},
				},
			},
		},
	}
}

func (c Newsboat) AddToCache(folders domain.Folders) error {
	tmpCachePath, db, closer, err := c.CreateNewTmpCache()
	defer closer()
	if err != nil {
		return err
	}

	var newsboatUrls []string

	// Add query to urls for listing all starred entries
	// Requires "prepopulate-query-feeds yes" in Newsboat config
	newsboatUrls = append(newsboatUrls, `"query:Starred:flags # \"s\""`)

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
				insertFeed,
				feed.ID, // id instead of url to disable manual refresh
				feed.Website,
				feed.Title,
			); err != nil {
				return err
			}

			log.Debug("Adding %d stories in feed %s", len(feed.Stories), feed.Title)
			for _, story := range feed.Stories {
				if _, err = db.Exec(
					insertStory,
					story.Hash, // our format is different, newsboat takes the <id> in <entry> if exists
					story.Title,
					story.Authors,
					story.Url,
					feed.ID, // id instead of url to disable manual refresh
					story.Timestamp.Unix(),
					story.Content,
					story.Unread,
					util.Cond(story.Starred, "s", ""),
				); err != nil {
					return err
				}
			}
		}
	}

	if err := util.MergeToFile(
		newsboatUrls,
		c.DataPath.GetFile("urls"),
		urlsSortFunc(),
	); err != nil {
		return err
	}

	return util.CopyFile(tmpCachePath, c.ReferenceDB(), c.UserDB())
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
