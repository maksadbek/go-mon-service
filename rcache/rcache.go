package rcache

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/logger"
	"github.com/garyburd/redigo/redis"
)

var (
	config conf.App        // config
	rc     ConcurrentRedis // redis client
)

// structure for fleet
type Fleet struct {
	Id      string           `json:"id"`           // unique id of fleet
	Update  map[string][]Pos `json:"update"`       // and its tracker's info
	LastReq int64            `json:"last_request"` // current unix time
}

// structure for fleet
type FleetTracker struct {
	Fleet    string
	Trackers []string
}

type Vehicle struct {
	Id                  int    `json:"id"`
	Fleet               int    `json:"fleet"`
	Imei                string `json:"imei"`
	Number              string `json:"number"`
	Tracker_type        string `json:"tracker_type"`
	Tracker_type_id     int    `json:"tracker_type_id"`
	Device_type_id      int    `json:"device_type_id"` // if this value is more than 0, then it has fuel sensor
	Name                string `json:"name"`
	Owner               string `json:"owner"`
	Active              string `json:"active"`
	Additional          string `json:"additional"`
	Customization       string `json:"customization"`
	Group_id            int    `json:"group_id"`
	Detector_fuel_id    int    `json:"detector_fuel_id"`
	Detector_motion_id  int    `json:"detector_motion_id"`
	Detector_dinamik_id int    `json:"detector_dinamik_id"`
	Pid                 int    `json:"pid"`
	Installed_sensor    string `json:"installed_sensor"`
	Detector_agro_id    int    `json:"detector_agro_id"`
	Car_health          string `json:"car_health"`
	Color               string `json:"color"`
	What_class          int    `json:"what_class"`
	ParamID             string `json:"a_param_id"`
}

func Initialize(c conf.App) (err error) {
	m := make(map[string]interface{})
	m["config"] = c
	logger.FuncLog("rcache.Initialize", "Initialize", m, nil)
	config = c
	// connect to redis
	err = rc.Start(c.DS.Redis.Host)
	if err != nil {
		logger.FuncLog("rcache.Initialize", "Unable to connect redis server", nil, err)
		return err
	}
	return
}

// GetTrackers can be used to get array of tracker of particular fleet
// start and stop are range values of list, default is 0,200, can be set from config
func GetTrackers(fleet string, start, stop int) (trackers []string, err error) {
	logger.FuncLog("rcache.GetTrackers", conf.InfoListOfTrackers, nil, nil)
	// get list of trackers ids from cache
	v, err := redis.Strings(rc.Do("SMEMBERS", config.DS.Redis.FPrefix+":"+fleet)) //
	if err != nil {
		logger.FuncLog("rcache.GetTrackers", conf.ErrGetListOfTrackers, nil, err)
		return
	}

	// range over list and append it to the slice
	for _, val := range v {
		trackers = append(trackers, val)
	}
	return
}

// PushRedis can be used to save fleet var into redis
func PushRedis(fleet Fleet) (err error) {
	logger.FuncLog("rcache.PushRedis", conf.InfoPushFleet, nil, nil)
	// get list of trackers ids from cache
	// range over map of Pos and push them
	for _, x := range fleet.Update {
		for _, pos := range x {
			jpos, err := json.Marshal(pos)
			if err != nil {
				logger.FuncLog("rcache.PushRedis", conf.ErrGetListOfTrackers, nil, err)
				return err
			}
			rc.Do("RPUSH", config.DS.Redis.TPrefix+":"+strconv.Itoa(pos.Id), jpos) // prefix can be set from conf
		}
	}
	return
}

// AddFleetTrackers can be used to save list of trackers to redis
func AddFleetTrackers(ftracker []FleetTracker) error {
	r, err := redis.Dial("tcp", config.DS.Redis.Host)
	if err != nil {
		logger.FuncLog("rcache.AddFleetTrackers", conf.ErrRedisConn, nil, err)
		return err
	}
	defer r.Close()
	// range over list
	for _, tracker := range ftracker {
		// range over tracker data
		for _, x := range tracker.Trackers {
			// add tracker to list
			r.Do("SADD", "fleet"+":"+tracker.Fleet, x)
		}
	}
	return nil
}

// CacheDefaults can be used to move all data in max_units table
// in mysql to redis
func CacheDefaults(trackers map[int]Vehicle) error {
	// create separate connection for caching
	r, err := redis.Dial("tcp", config.DS.Redis.Host)
	if err != nil {
		logger.FuncLog("rcache.AddFleetTrackers", conf.ErrRedisConn, nil, err)
		return err
	}
	defer r.Close()
	logger.FuncLog("rcache.CacheDefaults", "Caching Defaults", nil, nil)
	// range over map of data
	for id, x := range trackers {
		st := reflect.ValueOf(x)
		hashName := "max_unit_" + strconv.Itoa(id)
		for i := 0; i < st.NumField(); i++ {
			valueField := st.Field(i)
			typeField := st.Type().Field(i)
			key := fmt.Sprintf("%v", valueField.Interface())
			value := typeField.Name
			r.Do("HSET", hashName, value, key)
		}
	}
	return nil
}
