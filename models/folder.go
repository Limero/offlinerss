package models

type Folder struct {
	Id    int
	Title string
	Feeds Feeds
}

type Folders []*Folder

func (folders Folders) AddFolder(newFolder *Folder) (newFolders Folders) {
	/*
		Add a folder to a struct of folders in alphabetized order
	*/

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
