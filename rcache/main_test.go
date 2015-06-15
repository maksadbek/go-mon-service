package rcache

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"bitbucket.org/maksadbek/go-mon-service/conf"
)

func TestMain(m *testing.M) {
	r := strings.NewReader(mockConf) // читает мок-данные из testdata.go
	app, err := conf.Read(r)         //
	if err != nil {
		panic(err)
	}

	err = rc.Start(app.DS.Redis.Host)
	if err != nil {
		panic(err)
	}

	for _, x := range FleetTest.Trackers {
		rc.Do("SADD", "fleet"+":"+FleetTest.FleetName, x)
	}

	// add mock user
	jusr, err := json.Marshal(testUsr[0])
	if err != nil {
		panic(err)
	}
	rc.Do(
		"SET",
		app.DS.Redis.UPrefix+":"+testUsr[0].Login,
		string(jusr),
	)

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
