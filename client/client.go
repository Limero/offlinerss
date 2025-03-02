package client

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/limero/offlinerss/helpers"
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
)

type Client struct {
	ClientName   models.ClientName
	DataPath     models.DataPath
	Config       models.ClientConfig
	DatabaseInfo models.DatabaseInfo
	Files        models.ClientFiles
}

func (c Client) Name() models.ClientName {
	return c.ClientName
}

func (c Client) ReferenceDB() string {
	return c.DataPath.GetReferenceDB()
}

func (c Client) GetDatabaseInfo() models.DatabaseInfo {
	return c.DatabaseInfo
}

func (c Client) UserDB() string {
	return c.DataPath.GetFile(c.DatabaseInfo.FileName)
}

func (c Client) GetChanges() (models.SyncToActions, error) {
	return getChangesFromSqlite(
		c.ReferenceDB(),
		c.UserDB(),
		c.GetDatabaseInfo(),
	)
}

func (c Client) CreateNewCache() error {
	tmpCachePath := helpers.NewTmpCachePath()
	defer os.Remove(tmpCachePath)

	log.Debug("Creating %s temporary cache", c.ClientName)
	db, err := sql.Open(SQLiteDriver, tmpCachePath)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Debug("Creating tables in %s new temporary cache", c.ClientName)
	if _, err = db.Exec(c.DatabaseInfo.DDL); err != nil {
		return err
	}

	if err := helpers.CopyFile(tmpCachePath, c.ReferenceDB(), c.UserDB()); err != nil {
		return err
	}

	return nil
}

func (c Client) GetDataPath() models.DataPath {
	return c.DataPath
}

func (c Client) SetDataPath(dataPath models.DataPath) {
	c.DataPath = dataPath
}

func (c Client) GetFiles() models.ClientFiles {
	return c.Files
}

// Not exposed in interface
func (c Client) CreateNewTmpCache() (string, *sql.DB, func(), error) {
	tmpCachePath := helpers.NewTmpCachePath()

	if err := helpers.CopyFile(c.ReferenceDB(), tmpCachePath); err != nil {
		return "", nil, func() {}, err
	}

	closer := func() {
		os.Remove(tmpCachePath)
	}

	db, err := sql.Open(SQLiteDriver, tmpCachePath)
	if err != nil {
		return "", nil, closer, err
	}

	closer = func() {
		db.Close()
		os.Remove(tmpCachePath)
	}

	// Mark all items as read and unstarred, as we might never mark them otherwise
	// Everything currently unread and starred should be included in the folders we are adding later
	dbInfo := c.GetDatabaseInfo()
	_, err = db.Exec(fmt.Sprintf(
		"UPDATE %s SET %s = '%s', %s = '%s'",
		dbInfo.StoriesTable,
		dbInfo.Unread.Column,
		dbInfo.Unread.Negative,
		dbInfo.Starred.Column,
		dbInfo.Starred.Negative,
	))

	return tmpCachePath, db, closer, err
}
