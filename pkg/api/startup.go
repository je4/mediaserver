package api

import (
	"github.com/je4/mediaserver/v2/pkg/config"
	"github.com/je4/mediaserver/v2/pkg/databasePG"
	"github.com/je4/mediaserver/v2/pkg/service"
	"github.com/op/go-logging"
	"sync"
)

func Startup(conf *config.ApiConfig, dbClient *databasePG.Client, wg *sync.WaitGroup, log *logging.Logger) (service.ShutdownService, error) {
	ctrl, err := NewController(conf, nil, dbClient)
	if err != nil {
		log.Panicf("cannot create ingest controller: %v", err)
	}
	wg.Add(1)
	ctrl.Start(wg)
	return ctrl, nil
}
