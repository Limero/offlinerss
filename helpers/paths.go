package helpers

import (
	"os"
)

func GetMasterCachePath(clientName string) (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	return cacheDir + "/offlinerss/" + clientName + "/mastercache.db", nil
}
