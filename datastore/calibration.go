package datastore

import (
	log "bitbucket.org/maksadbek/go-mon-service/logger"
	"bitbucket.org/maksadbek/go-mon-service/rcache"
)

// LoadCalibres can be used to load calibration data form mysql
// and put it into storage
func LoadCalibres() error {
	log.FuncLog("datastore.LoadCalibres", "Loading calibration data", nil, nil)
	rcache.Calibres = make(map[int][]rcache.Calibration)
	rows, err := db.Query(queries["getCalibres"])
	if err != nil {
		return err
	}
	for rows.Next() {
		var c rcache.Calibration
		rows.Scan(
			&c.ID,
			&c.FleetID,
			&c.Litre,
			&c.Volt,
		)
		rcache.Calibres[c.ID] = append(rcache.Calibres[c.ID], c)
	}

	return nil
}
