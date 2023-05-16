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
	clientPath string,
	masterPath string,
	table string,
	idName string,
	unreadName string,
	unreadValueTrue string,
	unreadValueFalse string,
	starredName string,
	starredValueTrue string,
	starredValueFalse string,
) ([]models.SyncToAction, error) {
	if _, err := os.Stat(masterPath); os.IsNotExist(err) {
		fmt.Printf("Master cache does not exist at %s, nothing to sync to server\n", masterPath)
		return nil, nil
	}
	if _, err := os.Stat(clientPath); os.IsNotExist(err) {
		fmt.Printf("Cache does not exist at %s, nothing to sync to server\n", clientPath)
		return nil, nil
	}

	// Make copy of master cache to use for the sqldiff hack
	tmpCachePath := NewTmpCachePath()
	defer os.Remove(tmpCachePath)
	if err := CopyFile(masterPath, tmpCachePath); err != nil {
		return nil, err
	}

	fmt.Printf("Open master cache %s\n", masterPath)
	db, err := sql.Open("sqlite3", tmpCachePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// A one query workaround to get id to also show up in sqldiff
	if _, err = db.Exec("UPDATE " + table + " SET " + idName + "=NULL"); err != nil {
		return nil, err
	}

	fmt.Printf("Comparing database %q with %q\n", masterPath, clientPath)
	diffs, err := sqldiff.Compare(tmpCachePath, clientPath)
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
