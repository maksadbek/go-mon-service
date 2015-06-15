package datastore

import (
	"errors"
	"sync"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	log "bitbucket.org/maksadbek/go-mon-service/logger"
)

var (
	calibres map[int][]Calibration
	mutex    sync.RWMutex
)

type Calibration struct {
	ID      int
	FleetID int
	Litre   int
	Volt    float32
}

// LoadCalibres can be used to load calibration data form mysql
// and put it into storage
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

// GetLitrage can be used to get litrage value that is proportional
// to the voltage value of particular tracker
func GetLitrage(id int, volt float32) (litre int, err error) {
	mutex.RLock()
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
			numer := (int(volt) - int(calibre.Volt)) * (c[i+1].Litre - calibre.Litre)
			denom := int(c[i+1].Volt) - int(calibre.Volt)
			litre = numer/denom + calibre.Litre
			break
		}
	}
	mutex.RUnlock()
	return litre, err
}
