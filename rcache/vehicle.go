package rcache

import (
	"fmt"
	"sync"
)

type Vehicle struct {
	Id                  int               `json:"id"`
	Fleet               int               `json:"fleet"`
	Imei                string            `json:"imei"`
	Number              string            `json:"number"`
	Tracker_type        string            `json:"tracker_type"`
	Tracker_type_id     int               `json:"tracker_type_id"`
	Device_type_id      int               `json:"device_type_id"` // if this value is more than 0, then it has fuel sensor
	Name                string            `json:"name"`
	Owner               string            `json:"owner"`
	Active              string            `json:"active"`
	Additional          map[string]string `json:"additional"`
	Customization       string            `json:"customization"`
	Group_id            int               `json:"group_id"`
	Detector_fuel_id    int               `json:"detector_fuel_id"`
	Detector_motion_id  int               `json:"detector_motion_id"`
	Detector_dinamik_id int               `json:"detector_dinamik_id"`
	Pid                 int               `json:"pid"`
	Installed_sensor    string            `json:"installed_sensor"`
	Detector_agro_id    int               `json:"detector_agro_id"`
	Car_health          string            `json:"car_health"`
	Color               string            `json:"color"`
	What_class          int               `json:"what_class"`
	ParamID             string            `json:"a_param_id"`
}

type Vehicles struct {
	Data map[string]Vehicle
	sync.RWMutex
}

var VehicleList Vehicles

func (g *Vehicles) Put(id string, v Vehicle) {
	if len(g.Data) == 0 {
		g.Data = make(map[string]Vehicle)
	}
	g.Lock()
	g.Data[id] = v
	g.Unlock()
}

func (g *Vehicles) Get(id string) (Vehicle, error) {
	v, ok := g.Data[id]
	if !ok {
		return v, fmt.Errorf("vehicle (id %s) not found", id)
	}
	return v, nil
}
