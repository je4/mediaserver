package databasePG

import (
	"database/sql"
	"emperror.dev/errors"
	"encoding/json"
	"fmt"
	"github.com/je4/mediaserver/v2/pkg/models"
)

const LOAD_COLLECTIONS_ALL = "SELECT " +
	"collectionid, " +
	"estateid, " +
	"name, " +
	"description, " +
	"signature_prefix, " +
	"storageid, " +
	"jwtkey, " +
	"secret, " +
	"public " +
	"FROM %s.collection"

func NewCollectionsDB(db *sql.DB, schema string) (*collectionsDB, error) {
	cdb := &collectionsDB{
		db:     db,
		schema: schema,
	}
	return cdb, nil
}

type collectionsDB struct {
	db     *sql.DB
	schema string
}

func (cdb *collectionsDB) LoadAll(colls *models.Collections) error {
	colls.Lock()
	defer colls.Unlock()
	colls.Clear()
	sqlStr := fmt.Sprintf(
		LOAD_COLLECTIONS_ALL, cdb.schema)
	rows, err := cdb.db.Query(sqlStr)
	if err != nil {
		return errors.Wrapf(err, "cannot execute '%s'", sqlStr)
	}
	defer rows.Close()
	description := sql.NullString{}
	signature_prefix := sql.NullString{}
	jwtkey := sql.NullString{}
	secret := sql.NullString{}
	public := sql.NullString{}
	for rows.Next() {
		coll := &models.Collection{
			CollectionID:    0,
			EstateID:        0,
			Name:            "",
			Description:     "",
			SignaturePrefix: "",
			StorageID:       0,
			JWTKey:          "",
			Secret:          "",
			Public:          nil,
		}
		if err := rows.Scan(
			&coll.CollectionID,
			&coll.EstateID,
			&coll.Name,
			&description,
			&signature_prefix,
			&coll.StorageID,
			&jwtkey,
			&secret,
			&public,
		); err != nil {
			return errors.Wrapf(err, "cannot fetch row in query '%s'", sqlStr)
		}
		coll.Description = description.String
		coll.SignaturePrefix = signature_prefix.String
		coll.JWTKey = jwtkey.String
		coll.Secret = secret.String
		if public.Valid {
			var x any
			if err := json.Unmarshal([]byte(public.String), &x); err != nil {
				return errors.Wrapf(err, "invalid json \n%s\n", public.String)
			}
			coll.Public = x
		}
		colls.Add(coll)
	}
	return nil
}

var (
	_ models.CollectionsDatabase = (*collectionsDB)(nil)
)
