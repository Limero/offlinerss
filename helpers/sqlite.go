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
	clientConfig models.ClientConfig,
	table string,
	idName string,
	unreadName string,
	unreadValueTrue string,
	unreadValueFalse string,
	starredName string,
	starredValueTrue string,
	starredValueFalse string,
) ([]models.SyncToAction, error) {
	masterCachePath, err := GetMasterCachePath(clientConfig.Type)
	if err != nil {
		return nil, err
	}

	if _, err = os.Stat(masterCachePath); os.IsNotExist(err) {
		fmt.Printf("Master cache does not exist at %s, nothing to sync to server\n", masterCachePath)
		return nil, nil
	}
	if _, err = os.Stat(clientConfig.Paths.Cache); os.IsNotExist(err) {
		fmt.Printf("Cache does not exist at %s, nothing to sync to server\n", clientConfig.Paths.Cache)
		return nil, nil
	}

	fmt.Printf("Open master cache %s\n", masterCachePath)
	db, err := sql.Open("sqlite3", masterCachePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// A one query workaround to get id to also show up in sqldiff
	if _, err = db.Exec("UPDATE " + table + " SET " + idName + "=''"); err != nil {
		return nil, err
	}

	fmt.Printf("Comparing database %q with %q\n", masterCachePath, clientConfig.Paths.Cache)
	diffs, err := sqldiff.Compare(masterCachePath, clientConfig.Paths.Cache)
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
