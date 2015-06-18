package rcache

import (
	"encoding/json"
	"errors"
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

// structure for Tracker info
type Pos struct {
	Id            int     `json:"id"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
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
	Fuel          float32 `json:"fuel"`
	FuelVal       int     `json:"fuel_val"`
	MuAdditional  string  `json:"mu_additional"`
	Customization string  `json:"customization"`
	Additional    string  `json:"additional"`
	Action        int     `json:"action"`
}

// structure for fleet
type Fleet struct {
	Id      string         `json:"id"`           // unique id of fleet
	Update  map[string]Pos `json:"update"`       // and its tracker's info
	LastReq int64          `json:"last_request"` // current unix time
}

// structure for user
type Usr struct {
	Login    string   // login of user
	Fleet    string   // user's fleet id
	Trackers []string // user's list of trackers
}

// structure for fleet
type FleetTracker struct {
	Fleet    string
	Trackers []string
}

type Vehicle struct {
	Id                  int    `json=id`
	Fleet               int    `json=fleet`
	Imei                string `json=imei`
	Number              string `json=number`
	Tracker_type        string `json=tracker_type`
	Tracker_type_id     int    `json=tracker_type_id`
	Device_type_id      int    `json=device_type_id` // if this value is more than 0, then it has fuel sensor
	Name                string `json=name`
	Owner               string `json=owner`
	Active              string `json=active`
	DateCreated         string `json=dateCreated`
	Additional          string `json=additional`
	Customization       string `json=customization`
	Motor               int    `json=motor`
	MotorKoef           []byte `json=motorKoefbyte`
	CarSort             int    `json=carSort`
	Group_id            int    `json=group_id`
	YearOfManufac       int    `json=yearOfManufac`
	MarkerTypeId        int    `json=markerTypeId`
	ObjectTypeId        int    `json=objectTypeId`
	ScheduleId          int    `json=scheduleId`
	Detector_fuel_id    int    `json=detector_fuel_id`
	Detector_motion_id  int    `json=detector_motion_id`
	Detector_dinamik_id int    `json=detector_dinamik_id`
	Pid                 int    `json=pid`
	Installed_sensor    string `json=installed_sensor`
	Detector_agro_id    int    `json=detector_agro_id`
	Car_health          string `json=car_health`
	Color               string `json=color`
	What_class          int    `json=what_class`
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

// GetPositions can be used to retrieve map of positions
func GetPositions(trackerId []string) (trackers map[string]Pos, err error) {
	trackers = make(map[string]Pos)
	logger.FuncLog("rcache.GetPositions", "", nil, nil)
	// range over ids of trackers
	for _, tracker := range trackerId {
		var pos Pos
		// get tracker data
		pBytes, err := rc.Do("LINDEX", config.DS.Redis.TPrefix+":"+tracker, -1) // tracker's name saved with prefix, can be set from conf
		if err != nil {
			logger.FuncLog("rcache.GetPositions", err.Error(), nil, err)
			return trackers, err
		}
		p := fmt.Sprintf("%s", pBytes) // get string value of interface
		// if the value is nil, then merge with default values from max_units
		if fmt.Sprintf("%v", pBytes) == "<nil>" {
			logger.FuncLog("rcache.GetPositions", conf.ErrNotInCache, nil, nil)
			pos.Latitude = config.Defaults.Lat
			pos.Longitude = config.Defaults.Lng
			pos.Direction = config.Defaults.Direction
			pos.Speed = config.Defaults.Speed
			pos.Sat = config.Defaults.Sat
			pos.Ignition = config.Defaults.Ignition
			pos.GsmSignal = config.Defaults.GsmSignal
			pos.Battery = config.Defaults.Battery
			pos.Seat = config.Defaults.Seat
			pos.BatteryLvl = config.Defaults.BatteryLvl
			pos.Fuel = float32(config.Defaults.Fuel)
			pos.FuelVal = config.Defaults.FuelVal
			pos.MuAdditional = config.Defaults.MuAdditional
			pos.Action = config.Defaults.Action
			pos.Time = config.Defaults.Time

			idInt, err := strconv.Atoi(tracker)
			if err != nil {
				logger.FuncLog("rcache.GetPositions", "", nil, err)
			}
			pos.Id = idInt
			hashName := "max_unit_" + tracker

			// set default owner's name
			rOwner, err := rc.Do("HGET", hashName, "Owner")
			if err != nil {
				logger.FuncLog("rcache.GetPositions", conf.ErrNotInCache, nil, err)
			}
			if fmt.Sprintf("%v", rOwner) == "<nil>" {
				continue
			} else {
				pos.Owner = fmt.Sprintf("%s", rOwner)
			}

			// set default phone number
			rNumber, err := rc.Do("HGET", hashName, "Number")
			if err != nil {
				logger.FuncLog("rcache.GetPositions", conf.ErrNotInCache, nil, err)
			}
			if fmt.Sprintf("%v", rNumber) == "<nil>" {
				pos.Number = ""
			} else {
				pos.Number = fmt.Sprintf("%s", rNumber)
			}

			// set default name
			rName, err := rc.Do("HGET", hashName, "Name")
			if err != nil {
				logger.FuncLog("rcache.GetPositions", conf.ErrNotInCache, nil, err)
			}
			if fmt.Sprintf("%v", rName) == "<nil>" {
				pos.Name = ""
			} else {
				pos.Name = fmt.Sprintf("%s", rName)
			}

			// set default customization values
			rCustom, err := rc.Do("HGET", hashName, "Customization")
			if err != nil {
				logger.FuncLog("rcache.GetPositions", conf.ErrNotInCache, nil, err)
			}
			if fmt.Sprintf("%v", rCustom) == "<nil>" {
				pos.Customization = ""
			} else {
				pos.Customization = fmt.Sprintf("%s", rCustom)
			}

			// set default additional values
			rAdditional, err := rc.Do("HGET", hashName, "Additional")
			if err != nil {
				logger.FuncLog("rcache.GetPositions", conf.ErrNotInCache, nil, err)
			}
			if fmt.Sprintf("%v", rAdditional) == "<nil>" {
				pos.Additional = ""
			} else {
				pos.Additional = fmt.Sprintf("%s", rAdditional)
			}
			trackers[tracker] = pos
		} else {
			err = json.Unmarshal([]byte(p), &pos)
			if err != nil {
				logger.FuncLog("rcache.GetPositions", conf.ErrNotInCache, nil, err)
				return trackers, err
			}

			err = pos.SetLitrage()
			if err != nil {
				return trackers, err
			}
			trackers[tracker] = pos
		}
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
	for k, x := range fleet.Update {
		jpos, err := json.Marshal(x)
		if err != nil {
			logger.FuncLog("rcache.PushRedis", conf.ErrGetListOfTrackers, nil, err)
			return err
		}
		rc.Do("RPUSH", config.DS.Redis.TPrefix+":"+k, jpos) // prefix can be set from conf
	}
	return
}

// GetPositionsByFleet can be used to tracker data by fleet id
func GetPositionsByFleet(fleetNum string, start, stop int) (Fleet, error) {
	logger.FuncLog("rcache.PushRedis", "", nil, nil)
	fleet := Fleet{}
	fleet.Id = fleetNum
	fleet.Update = make(map[string]Pos)
	// get trackers of current fleet
	trackers, err := GetTrackers(fleetNum, start, stop)
	if err != nil {
		logger.FuncLog("rcache.GetPositionsByFleet", conf.ErrGetListOfTrackers, nil, err)
		return fleet, err
	}

	fleet.Update, err = GetPositions(trackers)
	if err != nil {
		logger.FuncLog("rcache.GetPositionsByFleet", conf.ErrGetListOfTrackers, nil, err)
	}
	return fleet, err
}

// UsrTrackers can be used to get info of user and list of its trackers
func UsrTrackers(name string) (Usr, error) {
	usr := Usr{}
	logger.FuncLog("rcache.UsrTrackers", "", nil, nil)
	// get user data
	userb, err := rc.Do("GET", config.DS.Redis.UPrefix+":"+name) // prefix can be set from conf
	if err != nil {
		logger.FuncLog("rcache.UsrTrackers", "", nil, err)
		return usr, err
	}
	// check whether it is nil, if nil then warn and finish
	if fmt.Sprintf("%v", userb) == "<nil>" {
		// prepare error message
		logger.FuncLog("rcache.UsrTrackers", conf.ErrNotInCache, nil, err)
		return usr, errors.New(config.ErrorMsg["NotExistInCache"].Msg)
	}
	err = json.Unmarshal([]byte(fmt.Sprintf("%s", userb)), &usr)
	if err != nil {
		logger.FuncLog("rcache.UsrTrackers", "", nil, err)
		return usr, err
	}
	return usr, nil
}

// SetUsrTrackers can be used to save user info in redis
func SetUsrTrackers(usr Usr) error {
	logger.FuncLog("rcache.SetUsrTrackers", "", nil, nil)
	jusr, err := json.Marshal(usr)
	if err != nil {
		logger.FuncLog("rcache.SetUsrTrackers", "", nil, err)
		return err
	}
	rc.Do(
		"SET",
		config.DS.Redis.UPrefix+":"+usr.Login,
		string(jusr),
	)
	return nil
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
