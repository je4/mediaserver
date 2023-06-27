package config

import (
	"emperror.dev/errors"
	"github.com/BurntSushi/toml"
	"github.com/je4/filesystem/v2/pkg/vfsrw"
)

type ApiConfig struct {
	Addr    string
	ExtAddr string
}

type Api struct {
	LogLevel   string
	LogFile    string
	API        ApiConfig    `toml:"API"`
	VFS        vfsrw.Config `toml:"VFS"`
	DatabasePG *DatabasePG
}

func LoadApiConfig(cfgData []byte) (*Api, error) {
	var config = &Api{LogLevel: "DEBUG"}
	if err := toml.Unmarshal(cfgData, config); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal toml")
	}
	return config, nil
}
