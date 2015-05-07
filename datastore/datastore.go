package datastore

import (
	"database/sql"
	"fmt"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	log "bitbucket.org/maksadbek/go-mon-service/logger"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Vehicle struct {
	id                  int
	fleet               int
	imei                string
	number              string
	tracker_type        string
	tracker_type_id     int
	device_type_id      int
	name                string
	owner               string
	active              string
	dateCreated         string
	additional          string
	customization       string
	motor               int
	motorKoef           []byte
	carSort             int
	group_id            int
	yearOfManufac       int
	markerTypeId        int
	objectTypeId        int
	scheduleId          int
	detector_fuel_id    int
	detector_motion_id  int
	detector_dinamik_id int
	pid                 int
	installed_sensor    string
	detector_agro_id    int
	car_health          string
	color               string
	what_class          int
}

func Initialize(c conf.Datastore) error {
	log.Log.WithFields(logrus.Fields{
		"package": "datastore",
		"config":  fmt.Sprintf("%+v", c),
	}).Info("Initialize")

	var err error
	db, err = sql.Open("mysql", c.Mysql.DSL)
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
	log.Log.WithFields(logrus.Fields{
		"package": "datastore",
		"fleet":   fleet,
	}).Info("GetTrackers")
	var pos map[int]Vehicle
	pos = make(map[int]Vehicle)
	query := `
			select 
				 id,
				 fleet,
				 imei,               
				 number,             
				 tracker_type,       
				 tracker_type_id,    
				 device_type_id,     
				 name,               
				 owner,              
				 active,             
				 additional,         
				 customization,      
				 group_id,           
				 detector_fuel_id,   
				 detector_motion_id, 
				 detector_dinamik_id,
				 pid,                
				 installed_sensor,   
				 detector_agro_id,   
				 car_health,         
				 color,              
				 what_class         
			from max_units
			where fleet = 202
	`
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
			&v.id,
			&v.fleet,
			&v.imei,
			&v.number,
			&v.tracker_type,
			&v.tracker_type_id,
			&v.device_type_id,
			&v.name,
			&v.owner,
			&v.active,
			&v.additional,
			&v.customization,
			&v.group_id,
			&v.detector_fuel_id,
			&v.detector_motion_id,
			&v.detector_dinamik_id,
			&v.pid,
			&v.installed_sensor,
			&v.detector_agro_id,
			&v.car_health,
			&v.color,
			&v.what_class,
		)
		pos[v.id] = v
	}
	log.Log.WithFields(logrus.Fields{
		"package":  "datastore",
		"postions": fmt.Sprintf("%+v", pos),
	}).Warn("GetTrackers")
	return pos, err
}

func CacheAllData() error {

}
