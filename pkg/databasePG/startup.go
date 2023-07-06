package databasePG

import (
	"database/sql"
	"emperror.dev/errors"
	"github.com/je4/mediaserver/v2/pkg/config"
	pb "github.com/je4/mediaserver/v2/pkg/protos"
	grpcutil "github.com/je4/utils/v2/pkg/grpc"
	_ "github.com/lib/pq"
	"github.com/op/go-logging"
	"google.golang.org/grpc"
	"os"
	"sync"
)

type databasePGShutdownService struct {
	db *sql.DB
	*grpc.Server
}

func (dss *databasePGShutdownService) Stop() {
	dss.Server.Stop()
	dss.db.Close()
}

func (dss *databasePGShutdownService) GracefulStop() {
	dss.Server.GracefulStop()
	dss.db.Close()
}

func Startup(conf *config.DatabasePG, wg *sync.WaitGroup, log *logging.Logger) (grpcutil.ShutdownService, error) {

	db, err := sql.Open("postgres", string(conf.Postgres.Connection))
	if err != nil {
		log.Panicf("cannot connect to database '%s': %v", conf.Postgres.Connection, err)
	}

	if err := db.Ping(); err != nil {
		log.Panicf("cannot ping database '%s': %v", conf.Postgres.Connection, err)
	}

	dbService, err := NewService(db, string(conf.Postgres.Schema))
	if err != nil {
		log.Panicf("cannot create database service: %v", err)
	}

	var serverCert, serverKey []byte
	if conf.ServerCert != "" {
		serverCert, err = os.ReadFile(string(conf.ServerCert))
		if err != nil {
			return nil, errors.Wrapf(err, "cannot read '%s'", conf.ServerCert)
		}
		serverKey, err = os.ReadFile(string(conf.ServerKey))
		if err != nil {
			return nil, errors.Wrapf(err, "cannot read '%s'", conf.ServerKey)
		}
	}

	shutter, err := grpcutil.Startup(string(conf.Addr), string(conf.Token), serverCert, serverKey, nil, func(srv *grpc.Server) {
		pb.RegisterDatabaseServer(srv, dbService)
	}, wg)
	if err != nil {
		log.Panicf("cannot start service at %s", conf.Addr)
	}

	return shutter, nil
}
