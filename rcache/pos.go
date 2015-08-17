package rcache

import (
	"encoding/json"
	"strconv"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/logger"
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
func GetPositions(trackerId []string) (map[string][]Pos, error) {
	trackers := make(map[string][]Pos)
	rc := pool.Get()
	defer rc.Close()
	// range over ids of trackers
	for _, id := range trackerId {
		var pos Pos
		var err error
		pos.Id, err = strconv.Atoi(id)
		if err != nil {
			logger.Log.Warn("GetPositions", err.Error())
		}
		// tracker's name saved with prefix, can be set from conf
		p, err := redis.String(rc.Do("LINDEX", config.DS.Redis.TPrefix+":"+id, -1))
		if err != nil {
			logger.Log.Error(err)
		}
		v, err := VehicleList.Get(id)
		if err != nil {
			logger.Log.Error(err)
			continue
		}
		// if the value is nil, then merge with default values from max_units
		if p == "" {
			// set default values
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

			pos.Owner = v.Owner
			pos.Number = v.Number
			pos.Name = v.Name
			pos.Customization = v.Customization
			j, _ := json.Marshal(v.Additional)
			pos.Additional = string(j)
		} else {
			err = json.Unmarshal([]byte(p), &pos)
			if err != nil {
				return trackers, err
			}
		}
		err = pos.SetLitrage(v.Device_type_id)
		if err != nil {
			return trackers, err
		}
		group, err := Grouplist.Get(strconv.Itoa(v.Group_id))
		if err != nil {
			group.Name = "all"
		}
		trackers[group.Name] = append(trackers[group.Name], pos)
	}
	return trackers, nil
}

// GetPositionsByFleet can be used to get tracker data by fleet id
func GetPositionsByFleet(fleetNum string, start, stop int) (Fleet, error) {
	fleet := Fleet{}
	fleet.Id = fleetNum
	fleet.Update = make(map[string][]Pos)
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
