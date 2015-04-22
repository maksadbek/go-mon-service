package conf

import (
	"github.com/BurntSushi/toml"
	"io"
)

type Datastore struct {
	Redis struct {
		Host    string // host of the redis server
		FPrefix string // prefix for fleet name in redis, i.e: in fleet_202, fleet is the prefix
	}
}

type Server struct {
	IP   string
	Port string
}

type App struct {
	DS  Datastore
	SRV Server
	Log Log
}

type Log struct {
	Path string
}

func Read(r io.Reader) (config App, err error) {
	_, err = toml.DecodeReader(r, &config)
	return config, err
}
