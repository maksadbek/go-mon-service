package datastore

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	log "bitbucket.org/maksadbek/go-mon-service/logger"
	"bitbucket.org/maksadbek/go-mon-service/rcache"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kovetskiy/go-php-serialize"
)

var (
	db       *sql.DB
	calibres map[int][]Calibration
	config   conf.App
)

type Calibration struct {
	ID      int
	FleetID int
	Litre   int
	Volt    float32
}

func Initialize(c conf.App) error {
	var err error
	config = c
	log.FuncLog("datastore.Initialize", "Initialization", nil, nil)
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

func GetTrackers(fleet string) (map[int]rcache.Vehicle, error) {
	queryFilter := " where fleet = " + fleet
	if fleet == "" {
		queryFilter = ""
	}
	log.Log.WithFields(logrus.Fields{
		"package": "datastore",
		"fleet":   fleet,
	}).Info("GetTrackers")
	var pos map[int]rcache.Vehicle = make(map[int]rcache.Vehicle)
	query := queries["getTrackers"] + queryFilter
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"package": "datastore",
			"Error":   err.Error(),
		}).Warn("GetTrackers")
		return pos, err
	}
	for rows.Next() {
		var v rcache.Vehicle
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
			&v.Additional,
			&v.Customization,
			&v.Group_id,
			&v.Detector_fuel_id,
			&v.Detector_motion_id,
			&v.Detector_dinamik_id,
			&v.Pid,
			&v.Installed_sensor,
			&v.Detector_agro_id,
			&v.Car_health,
			&v.Color,
			&v.What_class,
		)
		pos[v.Id] = v
	}
	log.Log.WithFields(logrus.Fields{
		"package": "datastore",
	}).Warn("GetTrackers")
	return pos, err
}

func UsrTrackers(name string) (usr rcache.Usr, err error) {
	log.Log.WithFields(logrus.Fields{
		"package": "datastore",
		"name":    name,
	}).Info("UsrTrackers")
	rows, err := db.Query(queries["usrTrackers"], name)
	defer rows.Close()
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"package": "datastore",
			"error":   err.Error(),
		}).Warn("UsrTrackers")
		return usr, err
	}

	for rows.Next() {
		var cars string
		rows.Scan(
			&usr.Login,
			&usr.Fleet,
			&cars,
		)
		if cars == "all" {
			usr.Trackers = append(usr.Trackers, "0")
			log.Log.WithFields(logrus.Fields{
				"package": "datastore",
				"user":    usr,
			}).Info("UsrTrackers")
			return usr, err
		}
		tr, err := phpserialize.Decode(cars)
		if err != nil {
			return usr, err
		}
		for _, x := range tr.(map[interface{}]interface{}) {
			usr.Trackers = append(usr.Trackers, fmt.Sprintf("%v", x))
		}
	}
	log.Log.WithFields(logrus.Fields{
		"package": "datastore",
		"user":    usr,
	}).Info("UsrTrackers")
	return
}

func CacheFleetTrackers() ([]rcache.FleetTracker, error) {
	log.Log.WithFields(logrus.Fields{
		"package": "datastore",
	}).Info("CacheFleetTrackers")
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

func LoadCalibres() error {
	log.FuncLog("datastore.LoadCalibres", "Loading calibration data", nil, nil)
	calibres = make(map[int][]Calibration)
	rows, err := db.Query(queries["getCalibres"])
	if err != nil {
		return err
	}
	for rows.Next() {
		var c Calibration
		rows.Scan(
			&c.ID,
			&c.FleetID,
			&c.Litre,
			&c.Volt,
		)
		calibres[c.ID] = append(calibres[c.ID], c)
	}

	return nil
}

func GetLitrage(id int, volt float32) (litre int, err error) {
	c := calibres[id]
	if c == nil {
		err = errors.New(conf.ErrNotInCache)
		return litre, err
	}
	for i, calibre := range c {
		if calibre.Volt == volt {
			litre = calibre.Litre
			return
		}
		if calibre.Volt < volt && c[i+1].Volt > volt {
			fmt.Printf("%v\n", calibre)
			fmt.Printf("%v\n", c[i+1])
			numer := (int(volt) - int(calibre.Volt)) * (c[i+1].Litre - calibre.Litre)
			denom := int(c[i+1].Volt) - int(calibre.Volt)
			fmt.Printf("numer is %d -> (%d - %d) * (%d - %d)\n", numer, int(volt), int(calibre.Volt), c[i+1].Litre, calibre.Litre)
			fmt.Printf("denom is %d -> (%d - %d) + %d\n", denom, int(c[i+1].Volt), int(calibre.Volt), calibre.Litre)
			litre = numer/denom + calibre.Litre
			break
		}
	}
	return litre, err
}
