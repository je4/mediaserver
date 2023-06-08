package databasePG

import (
	"context"
	"database/sql"
	"emperror.dev/errors"
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

	dbCollections, err := NewCollectionsDB(db, schema)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create collections db")
	}

	srv.collections, err = models.NewCollections(dbCollections)
	if err != nil {
		return nil, errors.Wrap(err, "cannot load collections")
	}

	dbStorages, err := NewStoragesDB(db, schema)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create storages db")
	}
	srv.storages, err = models.NewStorages(dbStorages)
	if err != nil {
		return nil, errors.Wrap(err, "cannot load storages")
	}

	return srv, nil
}

var _ pb.DatabaseServer = (*Service)(nil)
