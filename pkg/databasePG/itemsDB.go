package databasePG

import (
	"database/sql"
	"emperror.dev/errors"
	"fmt"
	"github.com/je4/mediaserver/v2/pkg/models"
	"github.com/lib/pq"
	"strings"
)

const LOAD_ITEM_SQL = "SELECT i.itemid, i.collectionid, i.signature, i.urn, i.type, i.subtype, i.objecttype, " +
	"i.parentid, i.mimetype, i.error, i.sha512, i.metadata, i.creation_type, i.last_modified, i.disabled, i.public, " +
	"i.public_actions, i.status FROM %s.item i, %s.collection c WHERE i.collectionid=c.collectionid AND c.name=? AND i.signature=?"

func NewItemsDB(db *sql.DB, schema string) (*itemsDB, error) {
	var err error
	cdb := &itemsDB{
		db:     db,
		schema: schema,
	}
	sqlStr := fmt.Sprintf(LOAD_ITEM_SQL, schema, schema)
	cdb.loadItemSQL, err = db.Prepare(sqlStr)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot prepare statement [%s]", sqlStr)
	}
	return cdb, nil
}

type itemsDB struct {
	db          *sql.DB
	schema      string
	loadItemSQL *sql.Stmt
}

func (cdb *itemsDB) LoadItem(collsig string) (*models.Item, error) {
	parts := strings.Split(collsig, "/")
	if len(parts) != 2 {
		return nil, errors.Errorf("invalid collection and signature '%s'", collsig)
	}
	var item = &models.Item{}
	coll := parts[0]
	sig := parts[1]
	row := cdb.loadItemSQL.QueryRow(coll, sig)
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
	_ models.ItemsDatabase = (*itemsDB)(nil)
)
