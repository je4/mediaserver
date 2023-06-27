package main

import (
	"emperror.dev/errors"
	"flag"
	"github.com/je4/mediaserver/v2/pkg/api"
	"github.com/je4/mediaserver/v2/pkg/config"
	"github.com/je4/mediaserver/v2/pkg/databasePG"
	lm "github.com/je4/utils/v2/pkg/logger"
	_ "github.com/lib/pq"
	"os"
	"sync"
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
	conf, err := config.LoadApiConfig(cfgData)
	if err != nil {
		panic(errors.Wrapf(err, "cannot unmarshal config toml data from '%s'", *configFile))
	}

	daLogger, lf := lm.CreateLogger("mediaserver-api", conf.LogFile, nil, conf.LogLevel, LOGFORMAT)
	defer lf.Close()

	dbClient, err := databasePG.NewClientPlain(string(conf.DatabasePG.Addr), "daToken")
	if err != nil {
		daLogger.Panicf("error creating database client for '%s'", string(conf.DatabasePG.Addr))
	}

	ctrl, err := api.NewController(&conf.API, nil, dbClient)
	if err != nil {
		daLogger.Panicf("cannot create ingest controller: %v", err)
	}
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	ctrl.Start(wg)

	wg.Wait()

	/*
		startFolder := "vfs:/digispace/ub-reprofiler/mets-container/bau1/2020"

		fs.WalkDir(vfs, startFolder, func(path string, d fs.DirEntry, err error) error {
			daLogger.Infof("path: %s", path)
			return nil
		})

		fp, err := vfs.Open("vfs:/digispace/ub-reprofiler/mets-container/bau1/2020/BAU_1_007097043_20190726T001152_master_ver1.zip/007097043/image/2316616.tif")
		if err != nil {
			daLogger.Panicf("cannot open tif in zip file")
		}
		defer fp.Close()
		fpw, err := writefs.Create(vfs, "vfs:/temp/test.tif")
		if err != nil {
			daLogger.Panicf("cannot create temp file")
		}
		defer fpw.Close()
		if _, err := io.Copy(fpw, fp); err != nil {
			daLogger.Panicf("error copying tif from zip")
		}

	*/

	/*
		dbService, err := database.NewService(db, conf.Postgres.Schema)
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

	*/
}
