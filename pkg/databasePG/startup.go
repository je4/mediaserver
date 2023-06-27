package databasePG

import (
	"database/sql"
	"fmt"
	"github.com/je4/mediaserver/v2/pkg/config"
	pb "github.com/je4/mediaserver/v2/pkg/protos"
	"github.com/je4/mediaserver/v2/pkg/service"
	grpcutil "github.com/je4/utils/v2/pkg/grpc"
	_ "github.com/lib/pq"
	"github.com/op/go-logging"
	"google.golang.org/grpc"
	"net"
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

func StartupPlain(conf *config.DatabasePG, wg *sync.WaitGroup, log *logging.Logger) (service.ShutdownService, error) {

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

	listener, err := net.Listen("tcp", string(conf.Addr))
	if err != nil {
		log.Panicf("cannot listen to tcp %s", conf.Addr)
	}

	var opts = []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcutil.JWTUnaryInterceptor),
		grpc.StreamInterceptor(grpcutil.JWTStreamInterceptor),
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterDatabaseServer(grpcServer, dbService)

	go func() {
		defer wg.Done()
		fmt.Printf("starting databasePG grpc server at %s\n", conf.Addr)
		if err := grpcServer.Serve(listener); err != nil {
			log.Errorf("error executing databasePG service: %v", err)
		}
		log.Info("databasePG service ended")
	}()
	return &databasePGShutdownService{
		db:     db,
		Server: grpcServer,
	}, nil
}

func StartupTLS(conf *config.DatabasePG, wg *sync.WaitGroup, log *logging.Logger) (service.ShutdownService, error) {

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

	listener, err := net.Listen("tcp", string(conf.Addr))
	if err != nil {
		log.Panicf("cannot listen to tcp %s", conf.Addr)
	}

	var opts = []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcutil.JWTUnaryInterceptor),
		grpc.StreamInterceptor(grpcutil.JWTStreamInterceptor),
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterDatabaseServer(grpcServer, dbService)

	go func() {
		defer wg.Done()
		fmt.Printf("starting databasePG grpc server at %s\n", conf.Addr)
		if err := grpcServer.Serve(listener); err != nil {
			log.Errorf("error executing databasePG service: %v", err)
		}
		log.Info("databasePG service ended")
	}()
	return &databasePGShutdownService{
		db:     db,
		Server: grpcServer,
	}, nil
}

func StartupConsul(conf *config.DatabasePG, wg *sync.WaitGroup, log *logging.Logger) (service.ShutdownService, error) {

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

	listener, err := net.Listen("tcp", string(conf.Addr))
	if err != nil {
		log.Panicf("cannot listen to tcp %s", conf.Addr)
	}

	var opts = []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcutil.JWTUnaryInterceptor),
		grpc.StreamInterceptor(grpcutil.JWTStreamInterceptor),
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterDatabaseServer(grpcServer, dbService)

	go func() {
		defer wg.Done()
		fmt.Printf("starting databasePG grpc server at %s\n", conf.Addr)
		if err := grpcServer.Serve(listener); err != nil {
			log.Errorf("error executing databasePG service: %v", err)
		}
		log.Info("databasePG service ended")
	}()
	return &databasePGShutdownService{
		db:     db,
		Server: grpcServer,
	}, nil
}
