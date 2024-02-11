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

// Default: ~/.config
func ConfigDir() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir != "" {
		return configDir
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, ".config")
}

// Default: ~/.local/share
func DataDir() string {
	configDir := os.Getenv("XDG_DATA_HOME")
	if configDir != "" {
		return configDir
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, ".local/share")
}
