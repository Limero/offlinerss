package main

import "github.com/limero/offlinerss/models"

func symlinkClientPaths(clients models.Clients) error {
	for _, client := range clients {
		for _, file := range client.GetFiles() {
			dataPath := client.GetDataPath()
			filePath := dataPath.GetFile(file.FileName)

			for _, targetPath := range file.TargetPaths {
				_ = filePath
				_ = targetPath

				// TODO
				// path exists and is correct symlink, skip
				// path exists and is incorrect symlink, unlink and relink
				// path exists and is file, rename and link
				// path doesn't exist, link
			}
		}
	}
	return nil
}
