package config

import (
	"emperror.dev/errors"
	"github.com/BurntSushi/toml"
	"github.com/je4/filesystem/v2/pkg/vfsrw"
)

type Ingest struct {
	LogLevel string
	LogFile  string
	Addr     string
	VFS      vfsrw.Config
}

func LoadIngestConfig(cfgData []byte) (*Ingest, error) {
	var config = &Ingest{LogLevel: "DEBUG", Addr: "localhost:8080"}
	if err := toml.Unmarshal(cfgData, config); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal toml")
	}
	return config, nil
}
