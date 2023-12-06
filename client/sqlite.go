package client

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
)

const (
	SQLiteDriver string = "sqlite3"
)

type dbRow struct {
	ID      string
	Unread  string
	Starred string
}

func getChangesFromSqlite(
	referenceDBPath string,
	userDBPath string,
	dbInfo models.DatabaseInfo,
) (models.SyncToActions, error) {
	if _, err := os.Stat(referenceDBPath); os.IsNotExist(err) {
		log.Debug("Reference database does not exist at %s, nothing to sync to server", referenceDBPath)
		return nil, nil
	}
	if _, err := os.Stat(userDBPath); os.IsNotExist(err) {
		log.Debug("User database does not exist at %s, nothing to sync to server", userDBPath)
		return nil, nil
	}
	refRows, err := getRowsFromDB(
		referenceDBPath,
		dbInfo.StoriesTable,
		dbInfo.StoriesIDColumn,
		dbInfo.Unread.Column,
		dbInfo.Starred.Column,
	)
	if err != nil {
		return nil, err
	}
	userRows, err := getRowsFromDB(
		userDBPath,
		dbInfo.StoriesTable,
		dbInfo.StoriesIDColumn,
		dbInfo.Unread.Column,
		dbInfo.Starred.Column,
	)
	if err != nil {
		return nil, err
	}

	var syncToActions models.SyncToActions
	for _, refRow := range refRows {
		for _, userRow := range userRows {
			if refRow.ID != userRow.ID {
				continue
			}

			if refRow.Unread != userRow.Unread {
				switch userRow.Unread {
				case dbInfo.Unread.Positive:
					syncToActions = append(syncToActions, models.SyncToAction{
						ID:     refRow.ID,
						Action: models.ActionStoryUnread,
					})
				case dbInfo.Unread.Negative:
					syncToActions = append(syncToActions, models.SyncToAction{
						ID:     refRow.ID,
						Action: models.ActionStoryRead,
					})
				}
			}

			if refRow.Starred != userRow.Starred {
				switch userRow.Starred {
				case dbInfo.Starred.Positive:
					syncToActions = append(syncToActions, models.SyncToAction{
						ID:     refRow.ID,
						Action: models.ActionStoryStarred,
					})
				case dbInfo.Starred.Negative:
					syncToActions = append(syncToActions, models.SyncToAction{
						ID:     refRow.ID,
						Action: models.ActionStoryUnstarred,
					})
				}
			}

			break
		}
	}

	return syncToActions, nil
}

func getRowsFromDB(dbPath, table, idName, unreadName, starredName string) ([]dbRow, error) {
	db, err := sql.Open(SQLiteDriver, dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(
		"SELECT %s, %s, %s FROM %s",
		idName,
		unreadName,
		starredName,
		table,
	)

	dbRows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer dbRows.Close()

	rows := make([]dbRow, 0)
	for dbRows.Next() {
		var r dbRow
		if err = dbRows.Scan(&r.ID, &r.Unread, &r.Starred); err != nil {
			return nil, err
		}
		rows = append(rows, r)
	}
	if err = dbRows.Err(); err != nil {
		return nil, err
	}

	return rows, nil
}
