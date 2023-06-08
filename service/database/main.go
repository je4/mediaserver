package main

import (
	"database/sql"
	"emperror.dev/errors"
	"flag"
	"fmt"
	"github.com/je4/mediaserver/v2/pkg/config"
	"github.com/je4/mediaserver/v2/pkg/databasePG"
	pb "github.com/je4/mediaserver/v2/pkg/protos"
	lm "github.com/je4/utils/v2/pkg/logger"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"net"
	"os"
)

const VERSION = "v1.0-beta.1"

const LOGFORMAT = `%{time:2006-01-02T15:04:05.000} %{shortpkg}::%{longfunc} [%{shortfile}] > %{level:.5s} - %{message}`

var configFile = flag.String("config", "./database.toml", "configuration file")

func main() {
	flag.Parse()

	cfgData, err := os.ReadFile(*configFile)
	if err != nil {
		panic(errors.Wrapf(err, "cannot read configuration from '%s'", *configFile))
	}
	conf, err := config.LoadDatabaseConfig(cfgData)
	if err != nil {
		panic(errors.Wrapf(err, "cannot unmarshal config toml data from '%s'", *configFile))
	}

	daLogger, lf := lm.CreateLogger("ocfl", conf.LogFile, nil, conf.LogLevel, LOGFORMAT)
	defer lf.Close()

	db, err := sql.Open("postgres", conf.Postgres.Connection)
	if err != nil {
		daLogger.Panicf("cannot connect to database '%s': %v", conf.Postgres.Connection, err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		daLogger.Panicf("cannot ping database '%s': %v", conf.Postgres.Connection, err)
	}

	dbService, err := databasePG.NewService(db, conf.Postgres.Schema)
	if err != nil {
		daLogger.Panicf("cannot create database service: %v", err)
	}

	listener, err := net.Listen("tcp", conf.Addr)
	if err != nil {
		daLogger.Panicf("cannot listen to tcp %s", conf.Addr)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterDatabaseServer(grpcServer, dbService)

	fmt.Printf("starting grpc server at %s", conf.Addr)
	grpcServer.Serve(listener)
}
