package datastore

var mockConf string = `[ds]
[ds.mysql]
    dsn = "root:toor@tcp(localhost:3306)/maxtrack"
    interval = 1
[ds.redis]
	host = ":6379"
	fprefix = "fleet"
    tprefix = "tracker"
[srv]
	port = ":1234"
[log]
	path = "info.log"
`
var FleetTest = struct {
	FleetName string
	Trackers  []string
}{
	"202",
	[]string{"106206", "107749", "107699"},
}
