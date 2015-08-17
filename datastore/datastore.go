package datastore

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	log "bitbucket.org/maksadbek/go-mon-service/logger"
	"bitbucket.org/maksadbek/go-mon-service/rcache"
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
	log.Log.Info("Datastore initialization")
	db, err = sql.Open("mysql", c.DS.Mysql.DSN)
	if err != nil {
		log.FuncLog("datastore.Initialize", "Initalization", nil, err)
		return err
	}
	err = LoadCalibres()
	if err != nil {
		log.FuncLog("datastore.Initialize", "Initalization", nil, err)
		return err
	}
	return nil
}

func CacheTrackers() error {
	log.Log.Info("Caching Default tracker data...")
	query := queries["getTrackers"]
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"package": "datastore",
			"Error":   err.Error(),
		}).Warn("GetTrackers")
		return err
	}
	for rows.Next() {
		var (
			v          rcache.Vehicle
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

		v.Additional = make(map[string]string)
		v.Additional = container
		v.ParamID = strconv.Itoa(int(paramId.Int64))
		v.Pid = int(pid.Int64)
		rcache.VehicleList.Put(strconv.Itoa(v.Id), v)
	}
	log.Log.Info("Succesfully cached default values")
	return nil
}

func UsrTrackers(name string) (usr rcache.Usr, err error) {
	var cars string
	log.Log.WithFields(logrus.Fields{
		"package": "datastore",
		"name":    name,
	}).Debug("UsrTrackers")
	err = db.QueryRow(queries["usrTrackers"], name).Scan(&usr.Login, &usr.Fleet, &cars)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"package": "datastore",
			"error":   err.Error(),
		}).Warn("UsrTrackers")
		return usr, err
	}

	if cars == "all" {
		usr.Trackers = append(usr.Trackers, "0")
		log.Log.WithFields(logrus.Fields{
			"package": "datastore",
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

func CacheFleetTrackers() ([]rcache.FleetTracker, error) {
	log.Log.WithFields(logrus.Fields{
		"package": "datastore",
	}).Debug("CacheFleetTrackers")
	var fleetTrackers []rcache.FleetTracker
	rows, err := db.Query(queries["fleetTrackers"])
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"package": "datastore",
			"error":   err.Error(),
		}).Warn("CacheFleetTrackers")
		return fleetTrackers, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			trackers rcache.FleetTracker
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
		log.FuncLog("datastore.CheckUser", "checking user", map[string]interface{}{"username": username, "hash": hash}, nil)
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
		group rcache.Group
		id    string
	)
	rows, err := db.Query(queries["trackerGroups"])
	defer rows.Close()
	if err != nil {
		log.FuncLog("datastore.LoadGroups", "checking user", nil, err)
		return err
	}
	for rows.Next() {
		rows.Scan(
			&id,
			&group.Name,
			&group.FleetID,
		)
		rcache.Grouplist.Put(id, group)
	}
	log.Log.Debug("Groups are loaded")
	return nil
}
