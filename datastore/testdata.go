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
