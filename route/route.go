package route

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/rcache"
)

func Initialize(c conf.App) error {
	err := rcache.Initialize(c)
	if err != nil {
		return err
	}
	return err
}
func GetPositionHandler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UTC().UnixNano())
	var testFleet rcache.Fleet = rcache.Fleet{
		Id: "202",
		Update: map[string]rcache.Pos{
			"106206": rcache.Pos{
				Id:            106206,
				Longitude:     69.145340,
				Latitude:      41.260006,
				Owner:         "Ozodbek",
				Number:        "01 S 775 JS",
				Name:          "Lacetti",
				Direction:     47,
				Speed:         10,
				Sat:           9,
				Time:          "2015-04-21 17:59:59",
				Ignition:      0,
				GsmSignal:     -1,
				Battery:       randInt(12000, 14000),
				Seat:          1000,
				BatteryLvl:    -1,
				Fuel:          randInt(40, 80),
				FuelVal:       randInt(10, 150),
				MuAdditional:  "",
				Customization: "a:1:{s:9:\"fillcolor\";s:7:\"#FF0000\";}",
				Additional:    "additional",
				Action:        1,
			},
			"107749": rcache.Pos{
				Id:            107749,
				Latitude:      41.260006,
				Longitude:     69.245811,
				Owner:         "Odil",
				Number:        "Acer Test",
				Name:          "Personal Acer",
				Direction:     243,
				Speed:         0,
				Sat:           99,
				Time:          "2015-03-26 13:29:06",
				Ignition:      0,
				GsmSignal:     -1,
				Battery:       randInt(12000, 14000),
				Seat:          1000,
				BatteryLvl:    -1,
				Fuel:          randInt(40, 80),
				FuelVal:       randInt(10, 150),
				MuAdditional:  "",
				Customization: "a:1:{s:9:\"fillcolor\";s:7:\"#FF0000\";}",
				Additional:    "additional",
				Action:        1,
			},
			"107699": rcache.Pos{
				Id:            107699,
				Longitude:     69.245926,
				Latitude:      41.293530,
				Owner:         "Odil",
				Number:        "01 048 QA",
				Name:          "PersonalAndroid",
				Direction:     225,
				Speed:         10,
				Sat:           99,
				Time:          "2015-04-01 18:00:13",
				Ignition:      0,
				GsmSignal:     -1,
				Battery:       randInt(12000, 14000),
				Seat:          1000,
				BatteryLvl:    -1,
				Fuel:          randInt(40, 80),
				FuelVal:       randInt(10, 150),
				MuAdditional:  "",
				Customization: "a:1:{s:9:\"fillcolor\";s:7:\"#FF0000\";}",
				Additional:    "",
				Action:        1,
			},
		},
	}
	pos, err := rcache.GetPositions("fleet_202", 0, 100)
	if err != nil {
		panic(err)
	}

	for key, x := range pos.Update {
		pos := testFleet.Update[key]
		pos.Latitude = x.Latitude + float32(0.00011)
		pos.Longitude = x.Longitude + float32(0.00011)
		testFleet.Update[key] = pos
	}
	jpos, err := json.Marshal(pos)
	if err != nil {
		panic(err)
	}

	err = rcache.PushRedis(testFleet)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(jpos))
}
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
