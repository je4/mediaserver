package databasePG

import (
	"database/sql"
	"emperror.dev/errors"
	"fmt"
	"github.com/je4/mediaserver/v2/pkg/models"
)

const LOAD_STORAGES_ALL = "SELECT " +
	"storageid, " +
	"name, " +
	"filebase, " +
	"datadir, " +
	"subitemdir, " +
	"tempdir " +
	"FROM %s.storage"

func NewStoragesDB(db *sql.DB, schema string) (*storagesDB, error) {
	cdb := &storagesDB{
		db:     db,
		schema: schema,
	}
	return cdb, nil
}

type storagesDB struct {
	db     *sql.DB
	schema string
}

func (cdb *storagesDB) LoadAll(stors *models.Storages) error {
	stors.Lock()
	defer stors.Unlock()
	stors.Clear()
	sqlStr := fmt.Sprintf(
		LOAD_STORAGES_ALL, cdb.schema)
	rows, err := cdb.db.Query(sqlStr)
	if err != nil {
		return errors.Wrapf(err, "cannot execute '%s'", sqlStr)
	}
	defer rows.Close()
	for rows.Next() {
		stor := &models.Storage{
			StorageID:  0,
			Name:       "",
			FileBase:   "",
			DataDir:    "",
			SubItemDir: "",
			TempDir:    "",
		}
		if err := rows.Scan(
			&stor.StorageID,
			&stor.Name,
			&stor.FileBase,
			&stor.DataDir,
			&stor.SubItemDir,
			&stor.TempDir,
		); err != nil {
			return errors.Wrapf(err, "cannot fetch row in query '%s'", sqlStr)
		}
		stors.Add(stor)
	}
	return nil
}

var (
	_ models.StoragesDatabase = (*storagesDB)(nil)
)
