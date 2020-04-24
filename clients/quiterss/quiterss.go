package quiterss

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
	err = conn.Exec("UPDATE news SET guid=''")
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
		if strings.Contains(row, " read=") {
			fmt.Printf("Change to read: %s\n", row)
			if !strings.Contains(row, "guid=") {
				if err != nil {
					return nil, errors.New("guid not found: " + row)
				}
			}
			hash := strings.Split(strings.Split(row, "guid='")[1], "'")[0]
			if strings.Contains(row, " read=2") {
				syncToActions = append(syncToActions, models.SyncToAction{
					Id:     hash,
					Action: models.ActionStoryRead,
				})
			} else if strings.Contains(row, " read=0") {
				syncToActions = append(syncToActions, models.SyncToAction{
					Id:     hash,
					Action: models.ActionStoryUnread,
				})
			} else {
				// should never be reached
				if err != nil {
					return nil, errors.New("read contains junk: " + row)
				}
			}
		}
	}
	return syncToActions, nil
}

func GenerateCache(folders []*models.Folder, clientConfig models.ClientConfig) error {
	tmpCachePath := fmt.Sprintf("%s/cache-%d.db", os.TempDir(), time.Now().UnixNano())
	defer os.Remove(tmpCachePath)

	fmt.Println("Creating QuiteRSS temporary cache")
	conn, err := sqlite3.Open(tmpCachePath)
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.BusyTimeout(5 * time.Second)

	fmt.Println("Creating tables in QuiteRSS new temporary cache")

	if err := conn.Exec(`
		CREATE TABLE feeds
			(
				id integer primary key,
				text varchar,
				title varchar,
				description varchar,
				xmlUrl varchar,
				htmlUrl varchar,
				language varchar,
				copyrights varchar,
				author_name varchar,
				author_email varchar,
				author_uri varchar,
				webMaster varchar,
				pubdate varchar,
				lastBuildDate varchar,
				category varchar,
				contributor varchar,
				generator varchar,
				docs varchar,
				cloud_domain varchar,
				cloud_port varchar,
				cloud_path varchar,
				cloud_procedure varchar,
				cloud_protocal varchar,
				ttl integer,
				skipHours varchar,
				skipDays varchar,
				image blob,
				unread integer,
				newCount integer,
				currentNews integer,
				label varchar,
				undeleteCount integer,
				tags varchar,
				hasChildren integer default 0,
				parentId integer default 0,
				rowToParent integer,
				updateIntervalEnable int,
				updateInterval int,
				updateIntervalType varchar,
				updateOnStartup int,
				displayOnStartup int,
				markReadAfterSecondsEnable int,
				markReadAfterSeconds int,
				markReadInNewspaper int,
				markDisplayedOnSwitchingFeed int,
				markDisplayedOnClosingTab int,
				markDisplayedOnMinimize int,
				layout text,
				filter text,
				groupBy int,
				displayNews int,
				displayEmbeddedImages integer default 1,
				loadTypes text,
				openLinkOnEmptyContent int,
				columns text,
				sort text,
				sortType int,
				maximumToKeep int,
				maximumToKeepEnable int,
				maximumAgeOfNews int,
				maximumAgoOfNewEnable int,
				deleteReadNews int,
				neverDeleteUnreadNews int,
				neverDeleteStarredNews int,
				neverDeleteLabeledNews int,
				status text,
				created text,
				updated text,
				lastDisplayed text,
				f_Expanded integer default 1,
				flags text,
				authentication integer default 0,
				duplicateNewsMode integer default 0,
				addSingleNewsAnyDateOn integer default 1,
				avoidedOldSingleNewsDateOn integer default 0,
				avoidedOldSingleNewsDate varchar,
				typeFeed integer default 0,
				showNotification integer default 0,
				disableUpdate integer default 0,
				javaScriptEnable integer default 1,
				layoutDirection integer default 0,
				SingleClickAction integer default 0,
				DoubleClickAction integer default 0,
				MiddleClickAction integer default 0
			);
		CREATE TABLE news
			(
				id integer primary key, feedId integer,
				guid varchar, guidislink varchar default 'true',
				description varchar, content varchar,
				title varchar, published varchar,
				modified varchar, received varchar,
				author_name varchar, author_uri varchar,
				author_email varchar, category varchar,
				label varchar, new integer default 1,
				read integer default 0, starred integer default 0,
				deleted integer default 0, attachment varchar,
				comments varchar, enclosure_length,
				enclosure_type, enclosure_url, source varchar,
				link_href varchar, link_enclosure varchar,
				link_related varchar, link_alternate varchar,
				contributor varchar, rights varchar,
				deleteDate varchar, feedParentId integer default 0
			)
	`); err != nil {
		return err
	}

	var unread int
	lastFeedId := 0 // This is required because folder/feed share same table and use ids

	fmt.Printf("Iterating over %d folders\n", len(folders))
	for _, folder := range folders {
		fmt.Printf("Add folder to database: %s\n", folder.Title)
		category := 0 // Category variable separate to lastFeedId to support feeds without a folder
		if folder.Title != "" {
			lastFeedId++
			category = lastFeedId
			if err := conn.Exec(
				"INSERT INTO feeds (id, text) VALUES (?, ?)",
				lastFeedId,
				folder.Title,
			); err != nil {
				return err
			}
		}

		fmt.Printf("Iterating over %d feeds in '%s' folder\n", len(folder.Feeds), folder.Title)
		for _, feed := range folder.Feeds {
			fmt.Printf("Add feed to database: %s\n", feed.Title)
			lastFeedId++
			if err := conn.Exec(
				"INSERT INTO feeds (id, text, title, xmlUrl, htmlUrl, unread, parentId) VALUES (?, ?, ?, ?, ?, ?, ?)",
				lastFeedId,
				feed.Title,
				feed.Title,
				feed.Url,
				feed.Website,
				feed.Unread,
				category,
			); err != nil {
				return err
			}

			fmt.Printf("Iterating over %d stories in feed %s\n", len(feed.Stories), feed.Title)
			for _, story := range feed.Stories {
				if story.Unread {
					unread = 0
				} else {
					unread = 2
				}

				fmt.Printf("\tAdd story to database: %s\n", story.Title)
				if err := conn.Exec(
					"INSERT INTO news (feedId, guid, description, title, published, read, link_href) VALUES (?, ?, ?, ?, ?, ?, ?)",
					lastFeedId,
					story.Hash,
					story.Content,
					story.Title,
					story.Timestamp,
					unread,
					story.Url,
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
