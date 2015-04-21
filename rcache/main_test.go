package rcache

import (
	"strings"
	"testing"

	"github.com/garyburd/redigo/redis"

	"bitbucket.org/maksadbek/go-mon-service/conf"
)

func TestMain(m *testing.M) {
	mockConf := `
	[ds]
		[ds.redis]
			host = ":6379"
	[srv]
		port = "1234"
	[log]
		path = "info.log"
	`
	r := strings.NewReader(mockConf)
	app, err := conf.Read(r)
	if err != nil {
		panic(err)
	}

	rc, err = redis.Dial("tcp", app.DS.Redis.Host)
	if err != nil {
		panic(err)
	}

	trackers := []string{"id123", "id456", "id789"}
	for _, x := range trackers {
		rc.Do("LPUSH", "fleet_202", x)
	}
	rc.Close()
	m.Run()
}
