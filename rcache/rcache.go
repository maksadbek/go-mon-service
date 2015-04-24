package rcache

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"github.com/garyburd/redigo/redis"
)

var (
	config conf.App   // config
	rc     redis.Conn // redis client
)

type Pos struct {
	Id        int
	Latitude  string
	Longitude string
	Time      string
}

type Fleet struct {
	Id     string
	Update map[string]Pos
}

func Initialize(c conf.App) (err error) {
	config = c
	rc, err = redis.Dial("tcp", c.DS.Redis.Host)
	if err != nil {
		return err
	}
	return
}

func GetTrackers(fleet string, start, stop int) (trackers []string, err error) {
	v, err := redis.Strings(rc.Do("LRANGE", fleet, start, stop))
	if err != nil {
		return
	}

	for _, val := range v {
		trackers = append(trackers, val)
	}
	return
}

func PushRedis(fleet Fleet) (err error) {
	for k, x := range fleet.Update {
		jpos, err := json.Marshal(x)
		if err != nil {
			return err
		}
		rc.Do("SADD", k, jpos)
	}
	return
}
func GetPositions(fleetNum string, start, stop int) (Fleet, error) {
	fleet := Fleet{}
	fleet.Id = fleetNum
	fleet.Update = make(map[string]Pos)
	trackers, err := GetTrackers(fleetNum, start, stop)
	if err != nil {
		return fleet, err
	}

	for _, v := range trackers {
		pBytes, err := redis.Values(rc.Do("SMEMBERS", v))
		if err != nil {
			return fleet, err
		}
		for _, x := range pBytes {
			p := fmt.Sprintf("%s", x)
			var pos Pos
			err = json.Unmarshal([]byte(p), &pos)
			if err != nil {
				return fleet, err
			}

			fleet.Update[v] = pos
		}
	}
	return fleet, err
}
