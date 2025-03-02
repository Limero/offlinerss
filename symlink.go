package main

import (
	"os"

	"github.com/limero/offlinerss/domain"
	"github.com/limero/offlinerss/helpers"
	"github.com/limero/offlinerss/log"
)

func symlinkClientPaths(clients domain.Clients) error {
	for _, client := range clients {
		for _, file := range client.GetFiles() {
			dataPath := client.GetDataPath()
			filePath := dataPath.GetFile(file.FileName)

			for _, targetPath := range file.TargetPaths {
				dest, _ := os.Readlink(targetPath)
				if dest != "" {
					if dest == filePath {
						log.Debug("Symlink from %q to %q is already correct", filePath, targetPath)
						continue
					}

					log.Warn("Removing incorrect symlink at %q", targetPath)
					if err := os.Remove(targetPath); err != nil {
						return err
					}
				} else if helpers.FileExists(targetPath) {
					log.Warn("Non-symlink found at target %q, renaming to .bak", targetPath)
					if err := os.Rename(targetPath, targetPath+".bak"); err != nil {
						return err
					}
				}

				if err := helpers.CreateParentDirs(targetPath); err != nil {
					return err
				}

				log.Info("Symlinking %q to %q", filePath, targetPath)
				if err := os.Symlink(filePath, targetPath); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
