package datastore

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	log "bitbucket.org/maksadbek/go-mon-service/logger"
	"bitbucket.org/maksadbek/go-mon-service/rcache"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Vehicle struct {
	Id                  int    `json=id`
	Fleet               int    `json=fleet`
	Imei                string `json=imei`
	Number              string `json=number`
	Tracker_type        string `json=tracker_type`
	Tracker_type_id     int    `json=tracker_type_id`
	Device_type_id      int    `json=device_type_id`
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

func Initialize(c conf.Datastore) error {
	log.Log.WithFields(logrus.Fields{
		"package": "datastore",
		"config":  fmt.Sprintf("%+v", c),
	}).Info("Initialize")

	var err error
	db, err = sql.Open("mysql", c.Mysql.DSN)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"package": "datastore",
			"Error":   err.Error(),
		}).Warn("Initialize")
		return err
	}
	return nil
}

func GetTrackers(fleet string) (map[int]Vehicle, error) {
	queryFilter := " where fleet = " + fleet
	if fleet == "" {
		queryFilter = ""
	}
	log.Log.WithFields(logrus.Fields{
		"package": "datastore",
		"fleet":   fleet,
	}).Info("GetTrackers")
	var pos map[int]Vehicle
	pos = make(map[int]Vehicle)
	query := ` select id, fleet, imei,               number,             tracker_type,       tracker_type_id,    device_type_id,     name,               owner,              active,             additional,         customization,      group_id,           detector_fuel_id,   detector_motion_id, detector_dinamik_id, pid,                installed_sensor,   detector_agro_id,   car_health,         color,              what_class         from max_units ` + queryFilter
	rows, err := db.Query(query)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"package": "datastore",
			"Error":   err.Error(),
		}).Warn("GetTrackers")
		return pos, err
	}
	for rows.Next() {
		var v Vehicle
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

func CacheData() error {
	trackers, err := GetTrackers("")
	if err != nil {
		return err
	}
	for id, x := range trackers {
		st := reflect.ValueOf(x)
		hashName := "max_unit_" + strconv.Itoa(id)
		for i := 0; i < st.NumField(); i++ {
			valueField := st.Field(i)
			typeField := st.Type().Field(i)
			key := fmt.Sprintf("%v", valueField.Interface())
			value := typeField.Name
			rcache.PutRawData(hashName, key, value)
		}
	}
	return err
}
