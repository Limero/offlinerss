package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func NewTmpCachePath() string {
	// Example: "/tmp/offlinerss-1684257550508849619.db"
	tmpFile := fmt.Sprintf("offlinerss-%d.db", time.Now().UnixNano())
	return filepath.Join(os.TempDir(), tmpFile)
}
