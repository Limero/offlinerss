package client

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/limero/offlinerss/domain"
	"github.com/limero/offlinerss/log"
	_ "github.com/mattn/go-sqlite3"
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
	dbInfo domain.DatabaseInfo,
) (domain.SyncToActions, error) {
	if _, err := os.Stat(referenceDBPath); os.IsNotExist(err) {
		log.Debug("Reference database does not exist at %s, nothing to sync to server", referenceDBPath)
		return domain.SyncToActions{}, nil
	}
	if _, err := os.Stat(userDBPath); os.IsNotExist(err) {
		log.Debug("User database does not exist at %s, nothing to sync to server", userDBPath)
		return domain.SyncToActions{}, nil
	}
	refRows, err := getRowsFromDB(
		referenceDBPath,
		dbInfo.StoriesTable,
		dbInfo.StoriesIDColumn,
		dbInfo.Unread.Column,
		dbInfo.Starred.Column,
	)
	if err != nil {
		return domain.SyncToActions{}, err
	}
	userRows, err := getRowsFromDB(
		userDBPath,
		dbInfo.StoriesTable,
		dbInfo.StoriesIDColumn,
		dbInfo.Unread.Column,
		dbInfo.Starred.Column,
	)
	if err != nil {
		return domain.SyncToActions{}, err
	}

	userRowsMap := make(map[string]dbRow, len(userRows))
	for _, row := range userRows {
		userRowsMap[row.ID] = row
	}

	var syncToActions domain.SyncToActions
	for _, refRow := range refRows {
		userRow, ok := userRowsMap[refRow.ID]
		if !ok {
			continue
		}

		if refRow.Unread != userRow.Unread {
			switch userRow.Unread {
			case dbInfo.Unread.Positive:
				syncToActions.Unread = append(syncToActions.Unread, refRow.ID)
			case dbInfo.Unread.Negative:
				syncToActions.Read = append(syncToActions.Read, refRow.ID)
			}
		}

		if refRow.Starred != userRow.Starred {
			switch userRow.Starred {
			case dbInfo.Starred.Positive:
				syncToActions.Starred = append(syncToActions.Starred, refRow.ID)
			case dbInfo.Starred.Negative:
				syncToActions.Unstarred = append(syncToActions.Unstarred, refRow.ID)
			}
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
