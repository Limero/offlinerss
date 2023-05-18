package test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/limero/offlinerss/client"
	"github.com/limero/offlinerss/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClients(t *testing.T) {
	tmpDir := os.TempDir()
	defer os.Remove(tmpDir)

	stories1 := []*models.Story{
		{
			Hash:   "123",
			Unread: true,
		},
		{
			Hash:   "321",
			Unread: true,
		},
	}
	stories2 := []*models.Story{
		{
			Hash:   "456",
			Unread: true,
		},
	}
	feeds := []*models.Feed{
		{
			Id:      123,
			Stories: stories1,
		},
	}
	folders := []*models.Folder{
		{
			Id:    123,
			Title: "Folder",
			Feeds: feeds,
		},
	}

	// TODO: Make client create this file
	os.Create(filepath.Join(tmpDir, "newsboat", "urls"))

	for _, tt := range []struct {
		client       models.Client
		storiesTable string
		updateQuery  string
	}{
		{
			client: client.Feedreader{
				DataPath: models.DataPath(filepath.Join(tmpDir, "feedreader")),
			},
			storiesTable: "articles",
			updateQuery:  "UPDATE articles SET unread = 8",
		},
		{
			client: client.Newsboat{
				DataPath: models.DataPath(filepath.Join(tmpDir, "newsboat")),
			},
			storiesTable: "rss_item",
			updateQuery:  "UPDATE rss_item SET unread = false",
		},
		{
			client: client.QuiteRSS{
				DataPath: models.DataPath(filepath.Join(tmpDir, "quiterss")),
			},
			storiesTable: "news",
			updateQuery:  "UPDATE news SET read = 2",
		},
	} {
		t.Run(tt.client.Name()+" create new cache", func(t *testing.T) {
			require.NoError(t, tt.client.CreateNewCache())
		})

		t.Run(tt.client.Name()+" get no changes on new cache", func(t *testing.T) {
			actions, err := tt.client.GetChanges()
			require.NoError(t, err)
			assert.Len(t, actions, 0)
		})

		t.Run(tt.client.Name()+" add folders to cache", func(t *testing.T) {
			require.NoError(t, tt.client.AddToCache(folders))

			db, err := sql.Open("sqlite3", tt.client.ReferenceDB())
			require.NoError(t, err)

			var count int
			err = db.QueryRow("SELECT COUNT(*) FROM " + tt.storiesTable).Scan(&count)
			require.NoError(t, err)
			assert.Len(t, stories1, count)
			require.NoError(t, db.Close())
		})

		t.Run(tt.client.Name()+" add same folders again to test idempotency", func(t *testing.T) {
			require.NoError(t, tt.client.AddToCache(folders))

			db, err := sql.Open("sqlite3", tt.client.ReferenceDB())
			require.NoError(t, err)

			var count int
			err = db.QueryRow("SELECT COUNT(*) FROM " + tt.storiesTable).Scan(&count)
			require.NoError(t, err)
			assert.Len(t, stories1, count)
			require.NoError(t, db.Close())
		})

		/*
			t.Run(tt.client.Name()+" add different folders test delta updates", func(t *testing.T) {
				folders[0].Feeds[0].Stories = stories2
				require.NoError(t, tt.client.AddToCache(folders))
				folders[0].Feeds[0].Stories = stories1

				db, err := sql.Open("sqlite3", tt.client.ReferenceDB())
				require.NoError(t, err)

				var count int
				err = db.QueryRow("SELECT COUNT(*) FROM " + tt.storiesTable).Scan(&count)
				require.NoError(t, err)
				assert.Equal(t, len(stories1)+len(stories2), count)
				require.NoError(t, db.Close())
			})
		*/
		_ = stories2

		t.Run(tt.client.Name()+" perform changes to user database", func(t *testing.T) {
			db, err := sql.Open("sqlite3", tt.client.UserDB())
			require.NoError(t, err)

			_, err = db.Exec(tt.updateQuery)
			require.NoError(t, err)
			require.NoError(t, db.Close())
		})

		t.Run(tt.client.Name()+" get changes performed", func(t *testing.T) {
			actions, err := tt.client.GetChanges()
			require.NoError(t, err)
			// Note: This is only events from stories2, as everything from stories1
			// was marked as read by the delta AddToCache call
			//assert.Len(t, actions, len(stories2))
			assert.Len(t, actions, len(stories1))
		})
	}
}
