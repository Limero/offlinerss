package models

type Folder struct {
	Id    int64
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

func (folders *Folders) GetOrCreateFolder(id int64, title string) *Folder {
	for _, folder := range *folders {
		if folder.Id == id {
			return folder
		}
	}

	newFolder := &Folder{
		Id:    id,
		Title: title,
		Feeds: Feeds{},
	}

	*folders = folders.AddFolder(newFolder)
	return newFolder
}

func (folders Folders) FindFeed(feedID int64) *Feed {
	for _, folder := range folders {
		for _, feed := range folder.Feeds {
			if feed.Id == feedID {
				return feed
			}
		}
	}
	return nil
}
