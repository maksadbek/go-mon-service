package rcache

import (
	"os"
	"strings"
	"testing"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"github.com/garyburd/redigo/redis"
)

func TestMain(m *testing.M) {
	r := strings.NewReader(mockConf) // читает мок-данные из testdata.go
	app, err := conf.Read(r)         // 
	if err != nil {
		panic(err)
	}

	rc, err = redis.Dial("tcp", app.DS.Redis.Host)
	if err != nil {
		panic(err)
	}

	for _, x := range FleetTest.Trackers {
		rc.Do("RPUSH", FleetTest.FleetName, x)
	}
	retCode := m.Run()

	// clean up messed redis test zone
	/*
		for _, x := range FleetTest.Trackers {
			rc.Do("LPOP", FleetTest.FleetName)
			rc.Do("LPOP", x)
		}
	*/

	os.Exit(retCode)
}
