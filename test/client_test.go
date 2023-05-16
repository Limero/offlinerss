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

	stories := []*models.Story{
		{
			Unread: true,
		},
	}
	feeds := []*models.Feed{
		{
			Id:      123,
			Stories: stories,
		},
	}
	folders := []*models.Folder{
		{
			Id:    123,
			Feeds: feeds,
		},
	}

	// TODO: Make client create this file
	os.Create(filepath.Join(tmpDir, "newsboat", "urls"))

	for _, tt := range []struct {
		client      models.Client
		updateQuery string
	}{
		{
			client: client.Feedreader{
				DataPath: models.DataPath(filepath.Join(tmpDir, "feedreader")),
			},
			updateQuery: "UPDATE articles SET unread = 8",
		},
		{
			client: client.Newsboat{
				DataPath: models.DataPath(filepath.Join(tmpDir, "newsboat")),
			},
			updateQuery: "UPDATE rss_item SET unread = false",
		},
		{
			client: client.QuiteRSS{
				DataPath: models.DataPath(filepath.Join(tmpDir, "quiterss")),
			},
			updateQuery: "UPDATE news SET read = 2",
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
		})

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
			assert.Len(t, actions, 1)
		})
	}
}
