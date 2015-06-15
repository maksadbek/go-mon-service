// @author: Maksadbek
// @email: a.maksadbek@gmail.com

package conf

import (
	"io"

	"github.com/BurntSushi/toml"
)

type ErrorStr struct {
	Msg string
}

type Auth struct {
	MACKey string
}
type Datastore struct {
	Mysql struct {
		DSN      string
		Interval int
	}
	Redis struct {
		Host     string
		FPrefix  string
		TPrefix  string
		UPrefix  string
		MUPrefix string
	}
}

const (
	ErrReq               = "request error"
	ErrNotInCache        = "not exist in cache"
	ErrGetListOfTrackers = "get list of tracker"
	ErrSetError          = "unable to set data into redis"
	ErrRedisConn         = "unable to connect redis server"
	ErrInvalidReq        = "invalid request data"
	ErrUnauthReq         = "unauthenticated request"
	InfoListOfTrackers   = "get list of trackers"
	InfoPushFleet        = "pushing fleet info"
)

type Server struct {
	IP   string
	Port string
}

type Defaults struct {
	Lat          float64
	Lng          float64
	Direction    int
	Speed        int
	Sat          int
	Ignition     int
	GsmSignal    int
	Battery      int
	Seat         int
	BatteryLvl   int
	Fuel         int
	FuelVal      int
	MuAdditional string
	Action       int
	Time         string
}

type App struct {
	DS       Datastore
	SRV      Server
	Log      Log
	ErrorMsg map[string]ErrorStr `toml:"errors"`
	Defaults Defaults
	Auth     Auth
}

type Log struct {
	Path string
}

func Read(r io.Reader) (config App, err error) {
	_, err = toml.DecodeReader(r, &config)
	return config, err
}
