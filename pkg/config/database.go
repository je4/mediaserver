package config

import (
	"emperror.dev/errors"
	"github.com/BurntSushi/toml"
	cfgutil "github.com/je4/utils/v2/pkg/config"
)

type Postgres struct {
	Connection cfgutil.EnvString `json:"connection"`
	Schema     cfgutil.EnvString `json:"schema"`
}

type Database struct {
	LogLevel cfgutil.EnvString
	LogFile  cfgutil.EnvString
	Addr     cfgutil.EnvString
	Postgres *Postgres
}

func LoadDatabaseConfig(cfgData []byte) (*Database, error) {
	var config = &Database{Addr: "localhost:1236", LogLevel: "DEBUG"}
	if err := toml.Unmarshal(cfgData, config); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal toml")
	}
	return config, nil
}
