package test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/limero/offlinerss/client"
	"github.com/limero/offlinerss/client/feedreader"
	"github.com/limero/offlinerss/client/newsboat"
	"github.com/limero/offlinerss/client/quiterss"
	"github.com/limero/offlinerss/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClients(t *testing.T) {
	for _, tt := range []struct {
		client        models.Client
		supportsDelta bool // TODO: Remove once all clients support delta updates
	}{
		{
			client:        feedreader.New(models.ClientConfig{}),
			supportsDelta: true,
		},
		{
			client:        newsboat.New(models.ClientConfig{}),
			supportsDelta: true,
		},
		{
			client:        quiterss.New(models.ClientConfig{}),
			supportsDelta: false,
		},
	} {
		stories1 := models.Stories{
			{
				Hash:    "123",
				Unread:  true,
				Starred: true,
			},
			{
				Hash:    "321",
				Unread:  true,
				Starred: false,
			},
		}
		stories2 := models.Stories{
			{
				Hash:    "456",
				Unread:  true,
				Starred: true,
			},
			{
				Hash:    "789",
				Unread:  false,
				Starred: false,
			},
		}
		feeds := models.Feeds{
			{
				ID:      123,
				Stories: stories1,
			},
		}
		folders := models.Folders{
			{
				ID:    123,
				Title: "Folder",
				Feeds: feeds,
			},
		}

		tt.client.SetDataPath(models.DataPath(t.TempDir()))

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
			expectDatabaseStories(t, tt.client, stories1)
		})

		t.Run(tt.client.Name()+" add same folders again to test idempotency", func(t *testing.T) {
			require.NoError(t, tt.client.AddToCache(folders))
			expectDatabaseStories(t, tt.client, stories1)
		})

		t.Run(tt.client.Name()+" add same folders again with different stories to test delta updates", func(t *testing.T) {
			folders[0].Feeds[0].Stories = stories2
			require.NoError(t, tt.client.AddToCache(folders))

			// Delta updates will mark all old stories as read/unstarred
			markedRead, markedUnstarred := 0, 0
			for i := range stories1 {
				if stories1[i].Unread {
					stories1[i].Unread = false
					markedRead++
				}
				if stories1[i].Starred {
					stories1[i].Starred = false
					markedUnstarred++
				}
			}
			// These asserts are to ensure this is being tested
			// instead of everything being read/unstarred already
			assert.Equal(t, 2, markedRead)
			assert.Equal(t, 1, markedUnstarred)

			expectedStories := stories2
			if tt.supportsDelta {
				expectedStories = append(stories1, stories2...)
			}
			expectDatabaseStories(t, tt.client, expectedStories)
		})

		t.Run(tt.client.Name()+" perform read change to user database", func(t *testing.T) {
			db, err := sql.Open(client.SQLiteDriver, tt.client.UserDB())
			require.NoError(t, err)

			dbInfo := tt.client.GetDatabaseInfo()
			res, err := db.Exec(fmt.Sprintf(
				"UPDATE %s SET %s = '%s' WHERE %s = %s",
				dbInfo.StoriesTable,
				dbInfo.Unread.Column,
				dbInfo.Unread.Negative,
				dbInfo.StoriesIDColumn,
				stories2[0].Hash,
			))
			require.NoError(t, err)

			rowsAffected, err := res.RowsAffected()
			require.NoError(t, err)
			assert.Equal(t, int64(1), rowsAffected)

			require.NoError(t, db.Close())
		})

		t.Run(tt.client.Name()+" perform unread change to user database", func(t *testing.T) {
			db, err := sql.Open(client.SQLiteDriver, tt.client.UserDB())
			require.NoError(t, err)

			dbInfo := tt.client.GetDatabaseInfo()
			res, err := db.Exec(fmt.Sprintf(
				"UPDATE %s SET %s = '%s' WHERE %s = %s",
				dbInfo.StoriesTable,
				dbInfo.Unread.Column,
				dbInfo.Unread.Positive,
				dbInfo.StoriesIDColumn,
				stories2[1].Hash,
			))
			require.NoError(t, err)

			rowsAffected, err := res.RowsAffected()
			require.NoError(t, err)
			assert.Equal(t, int64(1), rowsAffected)

			require.NoError(t, db.Close())
		})

		t.Run(tt.client.Name()+" perform unstarred change to user database", func(t *testing.T) {
			db, err := sql.Open(client.SQLiteDriver, tt.client.UserDB())
			require.NoError(t, err)

			dbInfo := tt.client.GetDatabaseInfo()

			res, err := db.Exec(fmt.Sprintf(
				"UPDATE %s SET %s = '%s' WHERE %s = %s",
				dbInfo.StoriesTable,
				dbInfo.Starred.Column,
				dbInfo.Starred.Negative,
				dbInfo.StoriesIDColumn,
				stories2[0].Hash,
			))
			require.NoError(t, err)

			rowsAffected, err := res.RowsAffected()
			require.NoError(t, err)
			assert.Equal(t, int64(1), rowsAffected)

			require.NoError(t, db.Close())
		})

		t.Run(tt.client.Name()+" perform starred change to user database", func(t *testing.T) {
			db, err := sql.Open(client.SQLiteDriver, tt.client.UserDB())
			require.NoError(t, err)

			dbInfo := tt.client.GetDatabaseInfo()

			res, err := db.Exec(fmt.Sprintf(
				"UPDATE %s SET %s = '%s' WHERE %s = %s",
				dbInfo.StoriesTable,
				dbInfo.Starred.Column,
				dbInfo.Starred.Positive,
				dbInfo.StoriesIDColumn,
				stories2[1].Hash,
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
			assert.Len(t, changes, 4)

			assert.Equal(t, stories2[0].Hash, changes[0].ID)
			assert.Equal(t, models.ActionStoryRead, changes[0].Action)

			assert.Equal(t, stories2[0].Hash, changes[1].ID)
			assert.Equal(t, models.ActionStoryUnstarred, changes[1].Action)

			assert.Equal(t, stories2[1].Hash, changes[2].ID)
			assert.Equal(t, models.ActionStoryUnread, changes[2].Action)

			assert.Equal(t, stories2[1].Hash, changes[3].ID)
			assert.Equal(t, models.ActionStoryStarred, changes[3].Action)
		})
	}
}

func expectDatabaseStories(t *testing.T, c models.Client, expectedStories models.Stories) {
	db, err := sql.Open(client.SQLiteDriver, c.ReferenceDB())
	require.NoError(t, err)

	dbInfo := c.GetDatabaseInfo()
	rows, err := db.Query(fmt.Sprintf(
		"SELECT %s, %s, %s FROM %s",
		dbInfo.StoriesIDColumn,
		dbInfo.Unread.Column,
		dbInfo.Starred.Column,
		dbInfo.StoriesTable,
	))
	require.NoError(t, err)
	defer rows.Close()

	count := 0
	for rows.Next() {
		var hash, unread, starred string
		require.NoError(t, rows.Scan(&hash, &unread, &starred))
		assert.Equal(t, expectedStories[count].Hash, hash)
		assert.Equal(t, expectedStories[count].Unread, unread == dbInfo.Unread.Positive)
		assert.Equal(t, expectedStories[count].Starred, starred == dbInfo.Starred.Positive)
		count++
	}
	assert.Len(t, expectedStories, count)

	require.NoError(t, db.Close())
}
