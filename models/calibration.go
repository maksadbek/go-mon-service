package models

import (
	"sort"

	log "github.com/Maksadbek/wherepo/logger"
	"github.com/Maksadbek/wherepo/cache"
)

type ByVoltage []cache.Calibration

func (a ByVoltage) Len() int           { return len(a) }
func (a ByVoltage) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByVoltage) Less(i, j int) bool { return a[i].Volt < a[j].Volt }

// LoadCalibres can be used to load calibration data form mysql
// and put it into storage
func LoadCalibres() error {
	log.FuncLog("models.LoadCalibres", "Loading calibration data", nil, nil)
	cache.Calibres = make(map[int][]cache.Calibration)
	cache.TopLitres = make(map[int]int)
	rows, err := db.Query(queries["getCalibres"])
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		var c cache.Calibration
		rows.Scan(
			&c.ID,
			&c.FleetID,
			&c.Litre,
			&c.Volt,
		)
		cache.Calibres[c.ID] = append(cache.Calibres[c.ID], c)
		sort.Sort(ByVoltage(cache.Calibres[c.ID]))
	}
	// load top litres
	rows, err = db.Query(queries["getTopLitres"])
	if err != nil {
		return err
	}
	for rows.Next() {
		var (
			id    int
			litre int
		)
		rows.Scan(
			&id,
			&litre,
		)
		cache.TopLitres[id] = litre
	}
	return nil
}
