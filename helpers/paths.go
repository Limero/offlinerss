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

// Defaults to ~/.config/file
func ConfigDir(file string) string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir != "" {
		return filepath.Join(configDir, file)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(filepath.Join(homeDir, ".config"), file)
}

// Defaults to ~/.local/share/file
func DataDir(file string) string {
	configDir := os.Getenv("XDG_DATA_HOME")
	if configDir != "" {
		return filepath.Join(configDir, file)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(filepath.Join(homeDir, ".local/share"), file)
}
