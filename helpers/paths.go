package helpers

import (
	"os"
)

func GetMasterCachePath(clientName string) string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	return cacheDir + "/offlinerss/" + clientName + "/mastercache.db"
}
