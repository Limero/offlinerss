package client

import (
	"database/sql"
	"os"

	"github.com/limero/offlinerss/helpers"
	"github.com/limero/offlinerss/log"
	"github.com/limero/offlinerss/models"
)

type Client struct {
	ClientName   string
	DataPath     models.DataPath
	Config       models.ClientConfig
	DatabaseInfo models.DatabaseInfo
}

func (c Client) Name() string {
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
	return helpers.GetChangesFromSqlite(
		c.ReferenceDB(),
		c.UserDB(),
		c.GetDatabaseInfo(),
	)
}

func (c Client) CreateNewCache() error {
	tmpCachePath := helpers.NewTmpCachePath()
	defer os.Remove(tmpCachePath)

	log.Debug("Creating %s temporary cache", c.ClientName)
	db, err := sql.Open("sqlite3", tmpCachePath)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Debug("Creating tables in %s new temporary cache", c.ClientName)
	if _, err = db.Exec(string(c.DatabaseInfo.DDL)); err != nil {
		return err
	}

	if err := helpers.CopyFile(tmpCachePath, c.ReferenceDB(), c.UserDB()); err != nil {
		return err
	}

	return nil
}

func (c Client) SetDataPath(dataPath models.DataPath) {
	c.DataPath = dataPath
}
