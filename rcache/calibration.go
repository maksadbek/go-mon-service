package rcache

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/logger"
)

var (
	Calibres  map[int][]Calibration
	TopLitres map[int]int
)

type Calibration struct {
	ID      int
	FleetID int
	Litre   int
	Volt    float32
}

// GetLitrage can be used to get litrage value that is proportional
// to the voltage value of particular tracker
func GetLitrage(id int, volt float32) (litre int, err error) {
	c := Calibres[id]
	if c == nil {
		err = errors.New(conf.ErrNotInCache)
		return litre, err
	}
	for i := 0; i < len(c)-1; i++ {
		calibre := c[i]
		if calibre.Volt == volt {
			litre = calibre.Litre
			return
		}
		if calibre.Volt < volt && c[i+1].Volt > volt {
			numer := (int(volt) - int(calibre.Volt)) * (c[i+1].Litre - calibre.Litre)
			denom := int(c[i+1].Volt) - int(calibre.Volt)
			litre = numer/denom + calibre.Litre
			break
		}
	}
	return litre, err
}

func (pos *Pos) SetLitrage(deviceTypeId int) error {
	id := strconv.Itoa(pos.Id)
	v, err := VehicleList.Get(id)
	if err != nil {
		logger.Log.Error(err)
	}

	additionals := make(map[string]float32)
	// validate id
	if id == "" {
		return errors.New("position id is nil")
	}

	if deviceTypeId > 0 {
		for _, x := range strings.Split(pos.Additional, ";") {
			m := strings.Split(x, ":")
			if len(m) == 2 {
				fuel, err := strconv.Atoi(m[1])
				if err != nil {
					logger.Log.Error(err)
					return err
				}
				additionals[m[0]] = float32(fuel)
			}
		}
		param := v.ParamID
		if err != nil {
			logger.Log.Error(err)
			return err
		}

		pos.FuelVal, err = GetLitrage(pos.Id, additionals[param])
		if err != nil {
			logger.Log.Error(err)
			return err
		}

		// set fuel value
		topLitre := TopLitres[pos.Id]

		// get percentage of FuelVal
		pos.Fuel = (100 * pos.FuelVal) / topLitre

		// set ignition
		// if error then set ignition = 0
		additionalsFromMU := make(map[string]string)
		add := v.Additional
		if err != nil {
			logger.Log.Error(err)
			pos.Ignition = 0
		}

		err = json.Unmarshal([]byte(add), &additionalsFromMU)
		if err != nil {
			logger.Log.Error(err)
			pos.Ignition = 0
		}

		ignitionIndex := additionalsFromMU["param_port_ignition"]
		ignition, ok := additionals[ignitionIndex]
		if !ok {
			pos.Ignition = 0
		}
		pos.Ignition = int(ignition)
	}
	return nil
}
