package helpers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/limero/offlinerss/models"
)

func SqlDiff(db1 string, db2 string) ([]string, error) {
	// Output SQL text that would transform DB1 into DB2
	fmt.Printf("Comparing database %s with %s\n", db1, db2)
	out, err := exec.Command("sqldiff", db1, db2).Output()
	if err != nil {
		return nil, fmt.Errorf("sqldiff failed: %w", err)
	}

	return strings.Split(string(out[:]), "\n"), nil
}

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

	if _, err := os.Stat(masterCachePath); os.IsNotExist(err) {
		fmt.Printf("Master cache does not exist at %s, nothing to sync to server\n", masterCachePath)
		return nil, nil
	}
	if _, err := os.Stat(clientConfig.Paths.Cache); os.IsNotExist(err) {
		fmt.Printf("Cache does not exist at %s, nothing to sync to server\n", clientConfig.Paths.Cache)
		return nil, nil
	}

	fmt.Printf("Open master cache %s\n", masterCachePath)
	conn, err := sqlite3.Open(masterCachePath)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	conn.BusyTimeout(5 * time.Second)

	// A one query workaround to get id to also show up in sqldiff
	err = conn.Exec("UPDATE " + table + " SET " + idName + "=''")
	if err != nil {
		return nil, err
	}

	diffs, err := SqlDiff(masterCachePath, clientConfig.Paths.Cache)
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
