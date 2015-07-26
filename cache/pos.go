package cache

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Maksadbek/wherepo/conf"
	"github.com/Maksadbek/wherepo/logger"
	"github.com/garyburd/redigo/redis"
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
	Fuel          int     `json:"fuel"`
	FuelVal       int     `json:"fuel_val"`
	MuAdditional  string  `json:"mu_additional"`
	Customization string  `json:"customization"`
	Additional    string  `json:"additional"`
	Action        int     `json:"action"`
}

// GetPositions can be used to retrieve map of positions
func GetPositions(trackerId []string) (trackers map[string][]Pos, err error) {
	trackers = make(map[string][]Pos)
	logger.FuncLog("cache.GetPositions", "", nil, nil)
	// range over ids of trackers
	for _, id := range trackerId {
		var pos Pos
		pos.Id, err = strconv.Atoi(id)
		if err != nil {
			logger.FuncLog("cache.GetPositions", "", nil, err)
		}
		// tracker's name saved with prefix, can be set from conf
		p, err := redis.String(rc.Do("LINDEX", config.DS.Redis.TPrefix+":"+id, -1))
		if err != nil {
			logger.Log.Error(err)
		}
		// get groupid of the tracker
		groupID, err := redis.String(rc.Do("HGET", "max_unit_"+id, "Group_id"))
		if err != nil {
			logger.Log.Error(err)
		}
		// if the value is nil, then merge with default values from max_units
		if p == "" {
			// set default values
			pos.SetPosDefaults()
		} else {
			err = json.Unmarshal([]byte(p), &pos)
			if err != nil {
				logger.FuncLog("cache.GetPositions", "Cannot unmarshal", nil, err)
				return trackers, err
			}
		}

		err = pos.SetLitrage()
		if err != nil {
			logger.Log.Error("here is it")
			return trackers, err
		}
		group, err := Grouplist.Get(groupID)
		if err != nil {
			logger.Log.Error(err)
			group.Name = "all"
		}
		trackers[group.Name] = append(trackers[group.Name], pos)
	}
	return
}

// GetPositionsByFleet can be used to get tracker data by fleet id
func GetPositionsByFleet(fleetNum string, start, stop int) (Fleet, error) {
	logger.FuncLog("cache.PushRedis", "", nil, nil)
	fleet := Fleet{}
	fleet.Id = fleetNum
	fleet.Update = make(map[string][]Pos)
	// get trackers of current fleet
	trackers, err := GetTrackers(fleetNum, start, stop)
	if err != nil {
		logger.FuncLog("cache.GetPositionsByFleet", conf.ErrGetListOfTrackers, nil, err)
		return fleet, err
	}

	fleet.Update, err = GetPositions(trackers)
	if err != nil {
		fmt.Println("error is in fleet.Update, err = GetPositions(trackers)")
		logger.FuncLog("cache.GetPositionsByFleet", conf.ErrGetListOfTrackers, nil, err)
	}
	return fleet, err
}

func (pos *Pos) SetPosDefaults() {
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
	pos.Fuel = config.Defaults.Fuel
	pos.FuelVal = config.Defaults.FuelVal
	pos.MuAdditional = config.Defaults.MuAdditional
	pos.Action = config.Defaults.Action
	pos.Time = config.Defaults.Time

	hashName := "max_unit_" + strconv.Itoa(pos.Id)
	// set default owner's name
	rOwner, err := redis.String(rc.Do("HGET", hashName, "Owner"))
	if err != nil {
		logger.FuncLog("cache.GetPositions", conf.ErrNotInCache, nil, err)
	}
	pos.Owner = rOwner

	// set default phone number
	rNumber, err := redis.String(rc.Do("HGET", hashName, "Number"))
	if err != nil {
		logger.FuncLog("cache.GetPositions", conf.ErrNotInCache, nil, err)
	}
	pos.Number = rNumber

	// set default name
	rName, err := redis.String(rc.Do("HGET", hashName, "Name"))
	if err != nil {
		logger.FuncLog("cache.GetPositions", conf.ErrNotInCache, nil, err)
	}
	pos.Name = rName

	// set default customization values
	rCustom, err := redis.String(rc.Do("HGET", hashName, "Customization"))
	if err != nil {
		logger.FuncLog("cache.GetPositions", conf.ErrNotInCache, nil, err)
	}
	pos.Customization = rCustom

	// set default additional values
	rAdditional, err := redis.String(rc.Do("HGET", hashName, "Additional"))
	if err != nil {
		logger.FuncLog("cache.GetPositions", conf.ErrNotInCache, nil, err)
	}
	pos.Additional = rAdditional
}
