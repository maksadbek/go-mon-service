package datastore

import (
	"sort"

	log "bitbucket.org/maksadbek/go-mon-service/logger"
	"bitbucket.org/maksadbek/go-mon-service/rcache"
)

type ByVoltage []rcache.Calibration

func (a ByVoltage) Len() int           { return len(a) }
func (a ByVoltage) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByVoltage) Less(i, j int) bool { return a[i].Volt < a[j].Volt }

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
		sort.Sort(ByVoltage(rcache.Calibres[c.ID]))
	}

	return nil
}
