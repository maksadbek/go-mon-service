package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Maksadbek/wherepo/conf"
	log "github.com/Maksadbek/wherepo/logger"
	"github.com/Maksadbek/wherepo/cache"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kovetskiy/go-php-serialize"
)

var (
	db     *sql.DB
	config conf.App
)

func Initialize(c conf.App) error {
	var err error
	config = c
	log.FuncLog("models.Initialize", "Initialization", nil, nil)
	db, err = sql.Open("mysql", c.DS.Mysql.DSN)
	if err != nil {
		log.FuncLog("models.Initialize", "Initalization", nil, err)
		return err
	}
	err = LoadCalibres()
	if err != nil {
		log.FuncLog("models.Initialize", "Initalization", nil, err)
		return err
	}
	return nil
}

func GetTrackers() (map[int]cache.Vehicle, error) {
	log.Log.WithFields(logrus.Fields{
		"package": "models",
	}).Debug("GetTrackers")
	var pos map[int]cache.Vehicle = make(map[int]cache.Vehicle)
	query := queries["getTrackers"]
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"package": "models",
			"Error":   err.Error(),
		}).Warn("GetTrackers")
		return pos, err
	}
	for rows.Next() {
		var (
			v          cache.Vehicle
			paramId    sql.NullInt64
			pid        sql.NullInt64
			additional string
		)
		rows.Scan(
			&v.Id,
			&v.Fleet,
			&v.Imei,
			&v.Number,
			&v.Tracker_type,
			&v.Tracker_type_id,
			&v.Device_type_id,
			&v.Name,
			&v.Owner,
			&v.Active,
			&additional,
			&v.Customization,
			&v.Group_id,
			&v.Detector_fuel_id,
			&v.Detector_motion_id,
			&v.Detector_dinamik_id,
			&pid,
			&v.Installed_sensor,
			&v.Detector_agro_id,
			&v.Car_health,
			&v.Color,
			&v.What_class,
			&paramId,
		)
		// unserialize additionals
		decodedMap, err := phpserialize.Decode(additional)
		container := make(map[string]string)
		if err != nil {
			log.Log.Error(err)
		}

		if decodedMap == nil {
			continue
		}
		// convert from map[interface{}]interface to map[string]string
		for key, val := range decodedMap.(map[interface{}]interface{}) {
			container[fmt.Sprintf("%v", key)] = fmt.Sprintf("%v", val)
		}

		jAdditionals, err := json.Marshal(container)
		if err != nil {
			panic(err)
		}

		v.Additional = string(jAdditionals)
		v.ParamID = strconv.Itoa(int(paramId.Int64))
		v.Pid = int(pid.Int64)
		pos[v.Id] = v
	}
	log.Log.WithFields(logrus.Fields{
		"package": "models",
	}).Warn("GetTrackers")
	return pos, err
}

func UsrTrackers(name string) (usr cache.Usr, err error) {
	var cars string
	log.Log.WithFields(logrus.Fields{
		"package": "models",
		"name":    name,
	}).Debug("UsrTrackers")
	err = db.QueryRow(queries["usrTrackers"], name).Scan(&usr.Login, &usr.Fleet, &cars)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"package": "models",
			"error":   err.Error(),
		}).Warn("UsrTrackers")
		return usr, err
	}

	if cars == "all" {
		usr.Trackers = append(usr.Trackers, "0")
		log.Log.WithFields(logrus.Fields{
			"package": "models",
			"user":    usr,
		}).Debug("UsrTrackers")
		return usr, err
	}

	// unserialize php-serialized array
	tr, err := phpserialize.Decode(cars)
	if err != nil {
		return usr, err
	}

	for _, x := range tr.(map[interface{}]interface{}) {
		usr.Trackers = append(usr.Trackers, fmt.Sprintf("%v", x))
	}
	return
}

func CacheFleetTrackers() ([]cache.FleetTracker, error) {
	log.Log.WithFields(logrus.Fields{
		"package": "models",
	}).Debug("CacheFleetTrackers")
	var fleetTrackers []cache.FleetTracker
	rows, err := db.Query(queries["fleetTrackers"])
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"package": "models",
			"error":   err.Error(),
		}).Warn("CacheFleetTrackers")
		return fleetTrackers, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			trackers cache.FleetTracker
			t        string
		)
		rows.Scan(
			&trackers.Fleet,
			&t,
		)
		trackers.Trackers = strings.Split(t, ",")
		fleetTrackers = append(fleetTrackers, trackers)
	}
	return fleetTrackers, nil
}

func CheckUser(username, hash string) (exists bool) {
	exists = false
	rows, err := db.Query(queries["checkUser"], username, hash)
	if err != nil {
		log.FuncLog("models.CheckUser", "checking user", map[string]interface{}{"username": username, "hash": hash}, nil)
	}
	defer rows.Close()

	for rows.Next() {
		exists = true
	}
	return
}

func LoadGroups() error {
	log.Log.Debug("Loading groups...")
	var (
		group cache.Group
		id    string
	)
	rows, err := db.Query(queries["trackerGroups"])
	defer rows.Close()
	if err != nil {
		log.FuncLog("models.LoadGroups", "checking user", nil, err)
		return err
	}
	for rows.Next() {
		rows.Scan(
			&id,
			&group.Name,
			&group.FleetID,
		)
		cache.Grouplist.Put(id, group)
	}
	log.Log.Debug("Groups are loaded")
	return nil
}
