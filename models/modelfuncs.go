package models

func AddFolderToFolders(folders []*Folder, newFolder *Folder) (newFolders []*Folder) {
	/*
		Add a folder to a struct of folders in alphabetized order
	*/

	if len(newFolder.Feeds) == 0 {
		// Skip empty folders
		return folders
	}

	newFolderAdded := false

	for _, folder := range folders {
		if !newFolderAdded && newFolder.Title < folder.Title {
			newFolders = append(newFolders, newFolder)
			newFolderAdded = true
		}

		newFolders = append(newFolders, folder)
	}

	if !newFolderAdded {
		return append(newFolders, newFolder)
	}

	return newFolders
}
