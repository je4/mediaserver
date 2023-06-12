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

type DatabasePG struct {
	Addr     cfgutil.EnvString
	Postgres *Postgres
}

type ServiceDatabasePG struct {
	LogLevel   cfgutil.EnvString
	LogFile    cfgutil.EnvString
	DatabasePG *DatabasePG
}

func LoadDatabasePGConfig(cfgData []byte) (*ServiceDatabasePG, error) {
	var config = &ServiceDatabasePG{LogLevel: "DEBUG"}
	if err := toml.Unmarshal(cfgData, config); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal toml")
	}
	return config, nil
}
