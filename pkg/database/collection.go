package database

import (
	"database/sql"
	"emperror.dev/errors"
	"encoding/json"
	"fmt"
	"sync"
)

func NewCollections(db *sql.DB, schema string) (*Collections, error) {
	var collections = &Collections{}
	return collections, collections.LoadAll(db, schema)
}

type Collections struct {
	sync.RWMutex
	collections map[string]*Collection
}

func (colls *Collections) Get(name string) (*Collection, error) {
	colls.RLock()
	defer colls.RUnlock()
	coll, ok := colls.collections[name]
	if !ok {
		return nil, errors.Wrapf(notFound, "collection '%s'", name)
	}
	return coll, nil
}

func (colls *Collections) LoadAll(db *sql.DB, schema string) error {
	colls.Lock()
	defer colls.Unlock()
	colls.collections = make(map[string]*Collection)
	sqlStr := fmt.Sprintf(
		"SELECT "+
			"collectionid, "+
			"estateid, "+
			"name, "+
			"description, "+
			"signature_prefix, "+
			"storageid, "+
			"jwtkey, "+
			"secret, "+
			"public "+
			"FROM %s.collection", schema)
	rows, err := db.Query(sqlStr)
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
		coll := &Collection{
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
		colls.collections[coll.Name] = coll
	}
	return nil
}

type Collection struct {
	CollectionID    int64
	EstateID        int64
	Name            string
	Description     string
	SignaturePrefix string
	StorageID       int64
	JWTKey          string
	Secret          string
	Public          any
}
