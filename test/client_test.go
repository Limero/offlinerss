package test

import (
	"database/sql"
	"fmt"
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
	tmpDir := filepath.Join(os.TempDir(), "offlinerss-clients")
	defer os.RemoveAll(tmpDir)

	stories1 := models.Stories{
		{
			Hash:   "123",
			Unread: true,
		},
		{
			Hash:   "321",
			Unread: true,
		},
	}
	stories2 := models.Stories{
		{
			Hash:    "456",
			Unread:  true,
			Starred: true,
		},
	}
	feeds := models.Feeds{
		{
			Id:      123,
			Stories: stories1,
		},
	}
	folders := models.Folders{
		{
			Id:    123,
			Title: "Folder",
			Feeds: feeds,
		},
	}

	for _, tt := range []struct {
		client        models.Client
		supportsDelta bool // TODO: Remove once all clients support delta updates
	}{
		{
			client: client.Feedreader{
				DataPath: models.DataPath(filepath.Join(tmpDir, "feedreader")),
			},
			supportsDelta: true,
		},
		{
			client: client.Newsboat{
				DataPath: models.DataPath(filepath.Join(tmpDir, "newsboat")),
			},
			supportsDelta: true,
		},
		{
			client: client.QuiteRSS{
				DataPath: models.DataPath(filepath.Join(tmpDir, "quiterss")),
			},
			supportsDelta: false,
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
			err = db.QueryRow("SELECT COUNT(*) FROM " + tt.client.GetDatabaseInfo().StoriesTable).Scan(&count)
			require.NoError(t, err)
			assert.Len(t, stories1, count)
			require.NoError(t, db.Close())
		})

		t.Run(tt.client.Name()+" add same folders again to test idempotency", func(t *testing.T) {
			require.NoError(t, tt.client.AddToCache(folders))

			db, err := sql.Open("sqlite3", tt.client.ReferenceDB())
			require.NoError(t, err)

			var count int
			err = db.QueryRow("SELECT COUNT(*) FROM " + tt.client.GetDatabaseInfo().StoriesTable).Scan(&count)
			require.NoError(t, err)
			assert.Len(t, stories1, count)
			require.NoError(t, db.Close())
		})

		t.Run(tt.client.Name()+" add same folders again with different stories to test delta updates", func(t *testing.T) {
			folders[0].Feeds[0].Stories = stories2
			require.NoError(t, tt.client.AddToCache(folders))
			folders[0].Feeds[0].Stories = stories1

			db, err := sql.Open("sqlite3", tt.client.ReferenceDB())
			require.NoError(t, err)

			var count int
			err = db.QueryRow("SELECT COUNT(*) FROM " + tt.client.GetDatabaseInfo().StoriesTable).Scan(&count)
			require.NoError(t, err)

			if tt.supportsDelta {
				assert.Equal(t, len(stories1)+len(stories2), count)
			} else {
				assert.Equal(t, len(stories2), count)
			}

			require.NoError(t, db.Close())
		})

		t.Run(tt.client.Name()+" perform read change to user database", func(t *testing.T) {
			db, err := sql.Open("sqlite3", tt.client.UserDB())
			require.NoError(t, err)

			dbInfo := tt.client.GetDatabaseInfo()
			res, err := db.Exec(fmt.Sprintf(
				"UPDATE %s SET %s = '%s' WHERE %s = %s",
				dbInfo.StoriesTable,
				dbInfo.Unread.Column,
				dbInfo.Unread.Negative,
				dbInfo.StoriesIdColumn,
				stories2[0].Hash,
			))
			require.NoError(t, err)

			rowsAffected, err := res.RowsAffected()
			require.NoError(t, err)
			assert.Equal(t, int64(1), rowsAffected)

			require.NoError(t, db.Close())
		})

		t.Run(tt.client.Name()+" perform unstarred change to user database", func(t *testing.T) {
			db, err := sql.Open("sqlite3", tt.client.UserDB())
			require.NoError(t, err)

			dbInfo := tt.client.GetDatabaseInfo()

			res, err := db.Exec(fmt.Sprintf(
				"UPDATE %s SET %s = '%s' WHERE %s = %s",
				dbInfo.StoriesTable,
				dbInfo.Starred.Column,
				dbInfo.Starred.Negative,
				dbInfo.StoriesIdColumn,
				stories2[0].Hash,
			))
			require.NoError(t, err)

			rowsAffected, err := res.RowsAffected()
			require.NoError(t, err)
			assert.Equal(t, int64(1), rowsAffected)

			require.NoError(t, db.Close())
		})

		t.Run(tt.client.Name()+" get changes performed", func(t *testing.T) {
			// Note: everything from stories1 was marked as read and unstarred
			// by the delta AddToCache call. So there won't be any changes to those
			changes, err := tt.client.GetChanges()
			require.NoError(t, err)
			assert.Len(t, changes, 2)

			assert.Equal(t, stories2[0].Hash, changes[0].Id)
			assert.Equal(t, models.ActionStoryRead, changes[0].Action)

			assert.Equal(t, stories2[0].Hash, changes[1].Id)
			assert.Equal(t, models.ActionStoryUnstarred, changes[1].Action)
		})
	}
}
