package client

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/limero/offlinerss/domain"
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/util"
)

type Client struct {
	ClientName   domain.ClientName
	DataPath     domain.DataPath
	Config       domain.ClientConfig
	DatabaseInfo domain.DatabaseInfo
	Files        domain.ClientFiles
}

func (c Client) Name() domain.ClientName {
	return c.ClientName
}

func (c Client) ReferenceDB() string {
	return c.DataPath.GetReferenceDB()
}

func (c Client) GetDatabaseInfo() domain.DatabaseInfo {
	return c.DatabaseInfo
}

func (c Client) UserDB() string {
	return c.DataPath.GetFile(c.DatabaseInfo.FileName)
}

func (c Client) GetChanges() (domain.SyncToActions, error) {
	return getChangesFromSqlite(
		c.ReferenceDB(),
		c.UserDB(),
		c.GetDatabaseInfo(),
	)
}

func (c Client) CreateNewCache() error {
	tmpCachePath := util.NewTmpCachePath()
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

	if err := util.CopyFile(tmpCachePath, c.ReferenceDB(), c.UserDB()); err != nil {
		return err
	}

	return nil
}

func (c Client) GetDataPath() domain.DataPath {
	return c.DataPath
}

func (c Client) SetDataPath(dataPath domain.DataPath) {
	c.DataPath = dataPath
}

func (c Client) GetFiles() domain.ClientFiles {
	return c.Files
}

// Not exposed in interface
func (c Client) CreateNewTmpCache() (string, *sql.DB, func(), error) {
	tmpCachePath := util.NewTmpCachePath()

	if err := util.CopyFile(c.ReferenceDB(), tmpCachePath); err != nil {
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
