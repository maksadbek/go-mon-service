package rcache

var FleetTest = struct {
	FleetName string
	Trackers  []string
}{
	"202",
	[]string{"106206", "107749", "107699"},
}
var mockConf string = `
[ds]
    [ds.redis]
		host = ":6379"
		fprefix = "fleet"
[srv]
    port = "1234"
[log]
    path = "info.log"
`
var testFleet Fleet = Fleet{
	Id: "202",
	Update: map[string]Pos{
		"106206": Pos{
			Id:            106206,
			Longitude:     69.145340,
			Owner:         "Ozodbek",
			Number:        "01 S 775 JS",
			Name:          "Lacetti",
			Latitude:      41.260006,
			Direction:     47,
			Speed:         10,
			Sat:           9,
			Time:          "2015-04-21 17:59:59",
			Ignition:      0,
			GsmSignal:     -1,
			Battery:       14293,
			Seat:          1000,
			BatteryLvl:    -1,
			Fuel:          0,
			FuelVal:       0,
			MuAdditional:  "",
			Customization: "",
			Additional:    "additional",
		},
		"107749": Pos{
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
			Battery:       5900,
			Seat:          1000,
			BatteryLvl:    -1,
			Fuel:          0,
			FuelVal:       0,
			MuAdditional:  "",
			Customization: "",
			Additional:    "additional",
		},
		"107699": Pos{
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
			Battery:       4500,
			Seat:          1000,
			BatteryLvl:    -1,
			Fuel:          0,
			FuelVal:       0,
			MuAdditional:  "",
			Customization: "a:1:{s:9:fillcolor;s:7:#993300;}",
			Additional:    "",
		},
	},
}
