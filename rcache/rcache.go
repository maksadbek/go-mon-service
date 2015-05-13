// @author: Maksadbek
// @email: a.maksadbek@gmail.com:
/*
   пакет для кеширования данных
*/

package rcache

import (
	"encoding/json"
	"fmt"
    "errors"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	log "bitbucket.org/maksadbek/go-mon-service/logger"
	"github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
)

var (
	config conf.App   // config
	rc     redis.Conn // redis client
)

// структура для позиции трекера
type Pos struct {
	Id            int     `json:"id"`
	Latitude      float32 `json:"latitude"`
	Longitude     float32 `json:"longitude"`
	Time          string  `json:"time"`
	Owner         string  `json:"owner"`
	Number        string  `json:"number"`
	Name          string  `json:"name"`
	Direction     int     `json:"direction"`
	Speed         int     `json:"speed"`
	Sat           int     `json:"sat"`
	Ignition      int     `json:"ignition"`
	GsmSignal     int     `json:"gsmsignal"`
	Battery       int     `json:"battery66"`
	Seat          int     `json:"seat"`
	BatteryLvl    int     `json:"batterylvl"`
	Fuel          int     `json:"fuel"`
	FuelVal       int     `json:"fuel_val"`
	MuAdditional  string  `json:"mu_additional"`
	Customization string  `json:"customization"`
	Additional    string  `json:"additional"`
	Action        int     `json:"action"`
}

// структура для флита
type Fleet struct {
	Id     string
	Update map[string]Pos
}

// функция для инициализации пакета
// оно должна вызыватся первые перед исползованием пакета
func Initialize(c conf.App) (err error) {
	log.Log.WithFields(logrus.Fields{
		"package": "rcache",
		"config":  fmt.Sprintf("%+v", c),
	}).Info("Initialization")

	config = c
	rc, err = redis.Dial("tcp", c.DS.Redis.Host)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"package":         "rcache",
			"Redis Dial host": c.DS.Redis.Host,
			"Error":           err.Error(),
		}).Fatal("Redis.Dial")
		return err
	}
	return
}

func GetPositions(trackerId []string) (trackers map[string]Pos, err error) {
	trackers = make(map[string]Pos)
	log.Log.WithFields(logrus.Fields{
		"package":  "rcache",
		"trackers": trackerId,
	}).Info("GetPositions")

	for _, tracker := range trackerId {
		fmt.Println(tracker)
		pBytes, err := rc.Do("LINDEX", config.DS.Redis.TPrefix+":"+tracker, -1)
		if err != nil {
			log.Log.WithFields(logrus.Fields{
				"package":       "rcache",
				"redis command": config.DS.Redis.FPrefix + ":" + tracker,
				"error":         err.Error(),
			}).Warn("GetPositions")
			return trackers, err
		}
		p := fmt.Sprintf("%s", pBytes)
        if fmt.Sprintf("%v", pBytes) == "<nil>" {
                return trackers, errors.New("nil value")
        }
		var pos Pos
		err = json.Unmarshal([]byte(p), &pos)
		if err != nil {
			log.Log.WithFields(logrus.Fields{
				"package": "rcache",
				"error":   err.Error(),
			}).Warn("GetPositions")
			return trackers, err
		}
		trackers[tracker] = pos
	}
	return
}

// функция используется для получения трекеров флита
func GetTrackers(fleet string, start, stop int) (trackers []string, err error) {
	log.Log.WithFields(logrus.Fields{
		"package": "rcache",
		"fleet":   fleet,
		"start":   start,
		"stop":    stop,
	}).Info("GetTrackers")

	v, err := redis.Strings(rc.Do("SMEMBERS", config.DS.Redis.FPrefix+":"+fleet))
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"package": "rcache",
			"Error":   err.Error(),
		}).Warn("GetTrackers")
		return
	}

	for _, val := range v {
		trackers = append(trackers, val)
	}
	return
}

// функция исползуется для вставления данный в редис
func PushRedis(fleet Fleet) (err error) {
	log.Log.WithFields(logrus.Fields{
		"package": "rcache",
	}).Info("PushRedis")
	for k, x := range fleet.Update {
		jpos, err := json.Marshal(x)
		if err != nil {
			log.Log.WithFields(logrus.Fields{
				"package": "rcache",
				"error":   err.Error(),
			}).Warn("PushRedis")
			return err
		}
		rc.Do("RPUSH", config.DS.Redis.TPrefix+":"+k, jpos)
	}
	return
}

func PutRawHash(hashName, field, data string) {
	log.Log.WithFields(logrus.Fields{
		"package": "rcache",
	}).Info("PushRawData")
	rc.Do("HSET", hashName, field, data)
	return
}

// исползуется для получения позиции трекеров флита
func GetPositionsByFleet(fleetNum string, start, stop int) (Fleet, error) {
	// log
	log.Log.WithFields(logrus.Fields{
		"package":     "rcache",
		"fleetNumber": fleetNum,
		"start":       start,
		"stop":        stop,
	}).Info("GetPositionsByFleet")
	fleet := Fleet{}
	fleet.Id = fleetNum
	fleet.Update = make(map[string]Pos)
	trackers, err := GetTrackers(fleetNum, start, stop)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"package": "rcache",
			"error":   err.Error(),
		}).Warn("GetPositionsByFleet")
		return fleet, err
	}

	for _, v := range trackers {
		pBytes, err := rc.Do("LINDEX", config.DS.Redis.TPrefix+":"+v, -1)
		if err != nil {
			log.Log.WithFields(logrus.Fields{
				"package": "rcache",
				"error":   err.Error(),
			}).Warn("GetPositionsByFleet")
			return fleet, err
		}
		p := fmt.Sprintf("%s", pBytes)
		var pos Pos
		err = json.Unmarshal([]byte(p), &pos)
		if err != nil {
			log.Log.WithFields(logrus.Fields{
				"package": "rcache",
				"error":   err.Error(),
			}).Warn("GetPositionsByFleet")
			return fleet, err
		}

		fleet.Update[v] = pos
	}
	return fleet, err
}

func FillPositions(p Pos) error {
	var err error
	return err
}
