package config

import (
	"emperror.dev/errors"
	"github.com/BurntSushi/toml"
	"os"
	"regexp"
	"strings"
)

type Postgres struct {
	Connection string `json:"connection"`
	Schema     string `json:"schema"`
}

type Database struct {
	LogLevel string
	LogFile  string
	Addr     string
	Postgres *Postgres
}

func LoadDatabaseConfig(cfgData []byte) (*Database, error) {
	var config = &Database{Addr: "localhost:1236", LogLevel: "DEBUG"}
	if err := toml.Unmarshal(cfgData, config); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal toml")
	}
	re := regexp.MustCompile(`%%([^%]+)%%`)
	matches := re.FindAllStringSubmatch(config.Postgres.Connection, -1)
	for _, match := range matches {
		data := os.Getenv(match[1])
		config.Postgres.Connection = strings.ReplaceAll(config.Postgres.Connection, match[0], data)
	}
	return config, nil
}
