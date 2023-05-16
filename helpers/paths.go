package helpers

import (
	"fmt"
	"os"
	"time"
)

func GetMasterCachePath(clientName string) string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	return cacheDir + "/offlinerss/" + clientName + "/mastercache.db"
}

func NewTmpCachePath() string {
	return fmt.Sprintf("%s/cache-%d.db", os.TempDir(), time.Now().UnixNano())
}
