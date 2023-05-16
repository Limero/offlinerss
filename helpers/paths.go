package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func GetClientFilePath(clientName string, file string) string {
	// Example: "~/.local/share/offlinerss/newsboat/cache.db"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(homeDir, ".local", "share", "offlinerss", clientName, file)
}

func GetMasterCachePath(clientName string) string {
	// Example: "~/.cache/offlinerss/newsboat/mastercache.db"
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(cacheDir, "offlinerss", clientName, "mastercache.db")
}

func NewTmpCachePath() string {
	// Example: "/tmp/offlinerss-1684257550508849619.db"
	tmpFile := fmt.Sprintf("offlinerss-%d.db", time.Now().UnixNano())
	return filepath.Join(os.TempDir(), tmpFile)
}
