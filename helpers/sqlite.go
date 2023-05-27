package helpers

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
)

type Row struct {
	Id      string
	Unread  string
	Starred string
}

func GetChangesFromSqlite(
	referenceDBPath string,
	userDBPath string,
	table string,
	idName string,
	unreadName string,
	unreadValueTrue string,
	unreadValueFalse string,
	starredName string,
	starredValueTrue string,
	starredValueFalse string,
) (models.SyncToActions, error) {
	if _, err := os.Stat(referenceDBPath); os.IsNotExist(err) {
		log.Debug("Reference database does not exist at %s, nothing to sync to server", referenceDBPath)
		return nil, nil
	}
	if _, err := os.Stat(userDBPath); os.IsNotExist(err) {
		log.Debug("User database does not exist at %s, nothing to sync to server", userDBPath)
		return nil, nil
	}

	refRows, err := getRowsFromDB(referenceDBPath, table, idName, unreadName, starredName)
	if err != nil {
		return nil, err
	}
	userRows, err := getRowsFromDB(userDBPath, table, idName, unreadName, starredName)
	if err != nil {
		return nil, err
	}

	var syncToActions models.SyncToActions
	for _, refRow := range refRows {
		for _, userRow := range userRows {
			if refRow.Id != userRow.Id {
				continue
			}

			if refRow.Unread != userRow.Unread {
				switch userRow.Unread {
				case unreadValueTrue:
					syncToActions = append(syncToActions, models.SyncToAction{
						Id:     refRow.Id,
						Action: models.ActionStoryUnread,
					})
				case unreadValueFalse:
					syncToActions = append(syncToActions, models.SyncToAction{
						Id:     refRow.Id,
						Action: models.ActionStoryRead,
					})
				}
			}

			if refRow.Starred != userRow.Starred {
				switch userRow.Starred {
				case starredValueTrue:
					syncToActions = append(syncToActions, models.SyncToAction{
						Id:     refRow.Id,
						Action: models.ActionStoryStarred,
					})
				case starredValueFalse:
					syncToActions = append(syncToActions, models.SyncToAction{
						Id:     refRow.Id,
						Action: models.ActionStoryUnstarred,
					})
				}
			}

			break
		}
	}

	return syncToActions, nil
}

func getRowsFromDB(dbPath, table, idName, unreadName, starredName string) ([]Row, error) {
	db, err := sql.Open("sqlite3", dbPath)
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

	rows := make([]Row, 0)
	for dbRows.Next() {
		var r Row
		if err = dbRows.Scan(&r.Id, &r.Unread, &r.Starred); err != nil {
			return nil, err
		}
		rows = append(rows, r)
	}
	if err = dbRows.Err(); err != nil {
		return nil, err
	}

	return rows, nil
}
