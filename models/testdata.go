package models

var mockConf string = `[ds]
[ds.mysql]
    dsn = "r23-go-mon:ErDFR8dqv44322www@tcp(maxtrack23:3306)/maxtrack"
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

var UserTest = struct {
	Username string
	Hash     string
}{
	"newmax",
	"f8cb56593dd08e04cd0f84d796b9cecd",
}
