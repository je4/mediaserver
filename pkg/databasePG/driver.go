package databasePG

import (
	"database/sql"
	"emperror.dev/errors"
	"encoding/json"
	"fmt"
	"github.com/je4/mediaserver/v2/pkg/models"
	"github.com/lib/pq"
)

const SQL_COLLECTIONS_LOADALL = "" +
	"SELECT " +
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

const SQL_STORAGES_LOADALL = "" +
	"SELECT " +
	"storageid, " +
	"name, " +
	"filebase, " +
	"datadir, " +
	"subitemdir, " +
	"tempdir " +
	"FROM %s.storage"

const SQL_LOADITEM = "" +
	"SELECT i.itemid, i.collectionid, i.signature, i.urn, i.type, i.subtype, i.objecttype, " +
	"i.parentid, i.mimetype, i.error, i.sha512, i.metadata, i.creation_date, i.last_modified, i.disabled, i.public, " +
	"i.public_actions, i.status " +
	"FROM %s.item i, %s.collection c " +
	"WHERE i.collectionid=c.collectionid AND c.name=$1 AND i.signature=$2"

const SQL_COLLECTION_NEWITEM = "" +
	"INSERT INTO %s.item (collectionid, signature, urn) " +
	"VALUES($1, $2, $3) returning itemid"

func NewDriver(db *sql.DB, schema string) (*driver, error) {
	var err error
	drv := &driver{
		db:     db,
		schema: schema,
	}
	sqlStr := fmt.Sprintf(SQL_LOADITEM, schema, schema)
	drv.loadItemSQL, err = db.Prepare(sqlStr)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot prepare statement [%s]", sqlStr)
	}

	return drv, nil
}

type driver struct {
	db          *sql.DB
	schema      string
	loadItemSQL *sql.Stmt
}

func (drv *driver) StoragesLoadAll(stores *models.Storages) error {
	stores.Lock()
	defer stores.Unlock()
	stores.Clear()
	sqlStr := fmt.Sprintf(
		SQL_STORAGES_LOADALL, drv.schema)
	rows, err := drv.db.Query(sqlStr)
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
		stores.Add(stor)
	}
	return nil
}

func (drv *driver) CollectionsLoadAll(cols *models.Collections) error {
	cols.Lock()
	defer cols.Unlock()
	cols.Clear()
	sqlStr := fmt.Sprintf(
		SQL_COLLECTIONS_LOADALL, drv.schema)
	rows, err := drv.db.Query(sqlStr)
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
		cols.Add(coll)
	}
	return nil
}

func (drv *driver) CollectionNewItem(collection models.Collection, signature, urn string) (*models.Item, error) {
	var params = []any{collection.CollectionID, signature, urn}
	row := drv.db.QueryRow(SQL_COLLECTION_NEWITEM, params...)
	var itemID int64
	if err := row.Scan(&itemID); err != nil {
		return nil, errors.Wrapf(err, "cannot scan result from '%s' %v", SQL_COLLECTION_NEWITEM, params)
	}
	return &models.Item{
		ItemID:       itemID,
		CollectionID: collection.CollectionID,
		Signature:    signature,
		Urn:          urn,
	}, nil
}

func (drv *driver) LoadItem(collection string, signature string) (*models.Item, error) {
	var item = &models.Item{}
	row := drv.loadItemSQL.QueryRow(collection, signature)
	var nType sql.NullString
	var nSubType sql.NullString
	var nParentID sql.NullInt64
	var nMimetype sql.NullString
	var nError sql.NullString
	var nSHA512 sql.NullString
	var nMetadata sql.NullString
	var nCreationDate pq.NullTime
	var nLastModified pq.NullTime
	var nDisabled sql.NullBool
	var nPublic sql.NullBool
	var nPublicActions pq.StringArray
	var nStatus sql.NullString

	if err := row.Scan(
		&item.ItemID,
		&item.CollectionID,
		&item.Signature,
		&item.Urn,
		&nType,
		&nSubType,
		&item.ObjectType,
		&nParentID,
		&nMimetype,
		&nError,
		&nSHA512,
		&nMetadata,
		&nCreationDate,
		&nLastModified,
		&nDisabled,
		&nPublic,
		&nPublicActions,
		&nStatus,
	); err != nil {
		return nil, errors.Wrap(err, "cannot scan result from load item query")
	}
	item.Type = nType.String
	item.SubType = nSubType.String
	item.ParentID = nParentID.Int64
	item.Mimetype = nMimetype.String
	item.Error = nError.String
	item.SHA512 = nSHA512.String
	item.Metadata = nMetadata.String
	item.CreationDate = nCreationDate.Time
	item.LastModified = nLastModified.Time
	item.Disbled = nDisabled.Bool
	item.Public = nPublic.Bool
	item.PublicActions = nPublicActions
	item.Status = nStatus.String

	return item, nil
}

var (
	_ models.ItemsDatabase       = (*driver)(nil)
	_ models.StoragesDatabase    = (*driver)(nil)
	_ models.CollectionsDatabase = (*driver)(nil)
)
