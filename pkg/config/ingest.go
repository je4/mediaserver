package config

import (
	"emperror.dev/errors"
	"github.com/BurntSushi/toml"
)

type Ingest struct {
	LogLevel string
	LogFile  string
}

func LoadIngestConfig(cfgData []byte) (*Ingest, error) {
	var config = &Ingest{LogLevel: "DEBUG"}
	if err := toml.Unmarshal(cfgData, config); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal toml")
	}
	return config, nil
}
