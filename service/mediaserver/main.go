package main

import (
	"context"
	"emperror.dev/errors"
	"flag"
	"github.com/je4/mediaserver/v2/pkg/api"
	"github.com/je4/mediaserver/v2/pkg/config"
	"github.com/je4/mediaserver/v2/pkg/databasePG"
	pb "github.com/je4/mediaserver/v2/pkg/protos"
	"github.com/je4/mediaserver/v2/pkg/service"
	lm "github.com/je4/utils/v2/pkg/logger"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const VERSION = "v1.0-beta.1"

const LOGFORMAT = `%{time:2006-01-02T15:04:05.000} %{shortpkg}::%{longfunc} [%{shortfile}] > %{level:.5s} - %{message}`

var configFile = flag.String("config", "./mediaserver.toml", "configuration file")

func main() {
	flag.Parse()

	cfgData, err := os.ReadFile(*configFile)
	if err != nil {
		panic(errors.Wrapf(err, "cannot read configuration from '%s'", *configFile))
	}
	conf, err := config.LoadMediaserverConfig(cfgData)
	if err != nil {
		panic(errors.Wrapf(err, "cannot unmarshal config toml data from '%s'", *configFile))
	}

	daLogger, lf := lm.CreateLogger("mediaserver", string(conf.LogFile), nil, string(conf.LogLevel), LOGFORMAT)
	defer lf.Close()

	var wg = &sync.WaitGroup{}

	var shutDownList = []service.ShutdownService{}

	sdl, err := databasePG.Startup(conf.DatabasePG, wg, daLogger)
	if err != nil {
		daLogger.Panicf("error starting databasPG: %v", err)
	}
	wg.Add(1)
	shutDownList = append(shutDownList, sdl)

	daLogger.Infof("sleeping 2sec")
	time.Sleep(2 * time.Second)

	dbClient, err := databasePG.NewClientPlain(string(conf.DatabasePG.Addr))
	if err != nil {
		daLogger.Panicf("error creating database client for '%s'", string(conf.DatabasePG.Addr))
	}

	if _, err := dbClient.Ping(context.Background(), &pb.Empty{}); err != nil {
		daLogger.Errorf("cannot ping database: %v", err)
	}
	sdl, err = api.Startup(conf.Api, dbClient, wg, daLogger)
	if err != nil {
		daLogger.Errorf("error starting databasPG: %v", err)
	}
	wg.Add(1)
	shutDownList = append(shutDownList, sdl)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-sigCh
		daLogger.Infof("got signal %v, attempting graceful shutdown", s)
		for _, sdl := range shutDownList {
			sdl.Stop()
		}
	}()

	wg.Wait()
}
