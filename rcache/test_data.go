package rcache

var FleetTest = struct {
	FleetName string
	Trackers  []string
}{
	"fleet_202",
	[]string{"id123", "id456", "id789"},
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
			Id:        106206,
			Longitude: "69.145340",
			Latitude:  "41.260006",
			Time:      "2015-04-21 17:59:59",
		},
		"107749": Pos{
			Id:        107749,
			Longitude: "69.245811",
			Latitude:  "41.293964",
			Time:      "2015-03-26 13:29:06",
		},
	},
}
