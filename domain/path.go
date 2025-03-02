package domain

import (
	"os"
	"path/filepath"
)

type DataPath string

func (d DataPath) GetFile(file string) string {
	return filepath.Join(string(d), file)
}

func (d DataPath) GetReferenceDB() string {
	return filepath.Join(string(d), "ref.db")
}

func GetClientDataPath(clientName ClientName) DataPath {
	// Example: "~/.local/share/offlinerss/newsboat/"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return DataPath(filepath.Join(homeDir, ".local", "share", "offlinerss", string(clientName)))
}
