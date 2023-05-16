package helpers

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/limero/go-sqldiff"
	"github.com/limero/offlinerss/models"
)

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
) ([]models.SyncToAction, error) {
	if _, err := os.Stat(referenceDBPath); os.IsNotExist(err) {
		fmt.Printf("Reference database does not exist at %s, nothing to sync to server\n", referenceDBPath)
		return nil, nil
	}
	if _, err := os.Stat(userDBPath); os.IsNotExist(err) {
		fmt.Printf("User database does not exist at %s, nothing to sync to server\n", userDBPath)
		return nil, nil
	}

	// Make copy of reference database to use for the sqldiff hack
	tmpCachePath := NewTmpCachePath()
	defer os.Remove(tmpCachePath)
	if err := CopyFile(referenceDBPath, tmpCachePath); err != nil {
		return nil, err
	}

	fmt.Printf("Opening reference database %s\n", referenceDBPath)
	db, err := sql.Open("sqlite3", tmpCachePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// A one query workaround to get id to also show up in sqldiff
	if _, err = db.Exec("UPDATE " + table + " SET " + idName + "=NULL"); err != nil {
		return nil, err
	}

	fmt.Printf("Comparing database %q with %q\n", referenceDBPath, userDBPath)
	diffs, err := sqldiff.Compare(tmpCachePath, userDBPath)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Iterating over %d database differences\n", len(diffs))
	var syncToActions []models.SyncToAction
	for _, row := range diffs {
		if strings.Contains(row, " "+unreadName+"=") {
			fmt.Printf("Change to %s: %s\n", unreadName, row)
			if !strings.Contains(row, idName+"=") {
				if err != nil {
					return nil, errors.New(idName + " not found: " + row)
				}
			}
			hash := strings.Split(strings.Split(row, idName+"='")[1], "'")[0]
			if strings.Contains(row, " "+unreadName+"="+unreadValueTrue) {
				syncToActions = append(syncToActions, models.SyncToAction{
					Id:     hash,
					Action: models.ActionStoryUnread,
				})
			} else if strings.Contains(row, " "+unreadName+"="+unreadValueFalse) {
				syncToActions = append(syncToActions, models.SyncToAction{
					Id:     hash,
					Action: models.ActionStoryRead,
				})
			}
		}
		if strings.Contains(row, " "+starredName+"=") {
			fmt.Printf("Change to %s: %s\n", starredName, row)
			if !strings.Contains(row, idName+"=") {
				if err != nil {
					return nil, errors.New(idName + " not found: " + row)
				}
			}
			hash := strings.Split(strings.Split(row, idName+"='")[1], "'")[0]

			if strings.Contains(row, " "+starredName+"="+starredValueTrue) {
				syncToActions = append(syncToActions, models.SyncToAction{
					Id:     hash,
					Action: models.ActionStoryStarred,
				})
			} else if strings.Contains(row, " "+starredName+"="+starredValueFalse) {
				syncToActions = append(syncToActions, models.SyncToAction{
					Id:     hash,
					Action: models.ActionStoryUnstarred,
				})
			}
		}
	}

	return syncToActions, nil
}
