package databasePG

import (
	"context"
	"database/sql"
	"emperror.dev/errors"
	"encoding/json"
	"fmt"
	"github.com/je4/mediaserver/v2/pkg/models"
	pb "github.com/je4/mediaserver/v2/pkg/protos"
)

type Service struct {
	db          *sql.DB
	schema      string
	collections *models.Collections
	storages    *models.Storages
	pb.UnimplementedDatabaseServer
}

func (srv *Service) Ping(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, srv.db.Ping()
}

func (srv *Service) GetCache(ctx context.Context, req *pb.CacheRequest) (*pb.CacheResult, error) {
	var result = &pb.CacheResult{
		Path:      "",
		Filesize:  0,
		Width:     nil,
		Height:    nil,
		Duration:  nil,
		MediaType: nil,
		Error:     nil,
	}

	coll, err := srv.collections.Get(req.Collection)
	if err != nil {
		var str = err.Error()
		result.Error = &str
		return result, nil
	}
	result.Error = &coll.Description

	return result, nil
}

func NewService(db *sql.DB, schema string) (*Service, error) {
	var err error
	srv := &Service{
		db:     db,
		schema: schema,
	}

	drv, err := NewDriver(db, schema)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create collections db")
	}

	srv.collections, err = models.NewCollections(drv)
	if err != nil {
		return nil, errors.Wrap(err, "cannot load collections")
	}
	srv.storages, err = models.NewStorages(drv)
	if err != nil {
		return nil, errors.Wrap(err, "cannot load storages")
	}

	return srv, nil
}

func (srv *Service) LoadAll(colls *models.Collections) error {
	colls.Lock()
	defer colls.Unlock()
	colls.Clear()
	sqlStr := fmt.Sprintf(
		LOAD_COLLECTIONS_ALL, srv.schema)
	rows, err := srv.db.Query(sqlStr)
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

var ()

var (
	_ pb.DatabaseServer          = (*Service)(nil)
	_ models.CollectionsDatabase = (*driver)(nil)
)
