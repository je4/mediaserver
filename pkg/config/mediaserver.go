package config

import (
	"emperror.dev/errors"
	"github.com/BurntSushi/toml"
	cfgutil "github.com/je4/utils/v2/pkg/config"
)

type ServiceMediaserver struct {
	LogLevel   cfgutil.EnvString
	LogFile    cfgutil.EnvString
	DatabasePG *DatabasePG
	Api        *ApiConfig
}

func LoadMediaserverConfig(cfgData []byte) (*ServiceMediaserver, error) {
	var config = &ServiceMediaserver{LogLevel: "DEBUG"}
	if err := toml.Unmarshal(cfgData, config); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal toml")
	}
	return config, nil
}
