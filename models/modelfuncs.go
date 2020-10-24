package models

func AddFolderToFolders(folders []*Folder, newFolder *Folder) (newFolders []*Folder) {
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

func AddFeedToFeeds(feeds []*Feed, newFeed *Feed) (newFeeds []*Feed) {
	/*
		Add a feed to a struct of feeds in alphabetized order
	*/

	newFeedAdded := false

	for _, feed := range feeds {
		if !newFeedAdded && newFeed.Title < feed.Title {
			newFeeds = append(newFeeds, newFeed)
			newFeedAdded = true
		}

		newFeeds = append(newFeeds, feed)
	}

	if !newFeedAdded {
		return append(newFeeds, newFeed)
	}

	return newFeeds
}
