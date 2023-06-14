package main

import (
	"emperror.dev/errors"
	"flag"
	"github.com/je4/filesystem/v2/pkg/vfsrw"
	"github.com/je4/filesystem/v2/pkg/writefs"
	"github.com/je4/mediaserver/v2/pkg/config"
	lm "github.com/je4/utils/v2/pkg/logger"
	_ "github.com/lib/pq"
	"io/fs"
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
	conf, err := config.LoadIngestConfig(cfgData)
	if err != nil {
		panic(errors.Wrapf(err, "cannot unmarshal config toml data from '%s'", *configFile))
	}

	daLogger, lf := lm.CreateLogger("ingest", string(conf.LogFile), nil, string(conf.LogLevel), LOGFORMAT)
	defer lf.Close()

	vfs, err := vfsrw.NewFS(conf.VFS, daLogger)
	if err != nil {
		daLogger.Panicf("cannot create vfs: %v", err)
	}
	listener, err := net.Listen("tcp", string(conf.Addr))
	if err != nil {
		daLogger.Panicf("cannot listen to tcp %s", conf.Addr)
	}

	if err := writefs.WriteFile(vfs, "vfs:/switch_ch/testarchive/data/hello.txt", []byte("hello")); err != nil {
		daLogger.Errorf("error writing file")
	}

	fs.WalkDir(vfs, "vfs:/switch_ch/", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			daLogger.Infof("[d] %s", path)
		} else {
			daLogger.Infof("[f] %s", path)
		}

		return nil
	})

	_ = vfs
	_ = listener

	/*
		var opts []grpc.ServerOption
		grpcServer := grpc.NewServer(opts...)
		pb.RegisterDatabaseServer(grpcServer, dbService)

		fmt.Printf("starting grpc server at %s", conf.Addr)
		grpcServer.Serve(listener)

	*/
}
