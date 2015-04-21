package rcache

import (
	"bitbucket.org/maksadbek/go-mon-service/conf"
	"github.com/garyburd/redigo/redis"
)

var (
	config conf.App   // config
	rc     redis.Conn // redis client
)

func Initialize(c conf.App) (err error) {
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
