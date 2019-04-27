package config

import (
	"github.com/BurntSushi/toml"
)

type Driver struct {
	// Binary  string              `toml:"binary"`
	Options map[string][]string `toml:"options"`
}

type Config struct {
	Driver *Driver
}

func Load(file string) (config *Config) {
	if _, err := toml.DecodeFile(file, &config); nil != err {
		panic(err)
	}
	return config
}
