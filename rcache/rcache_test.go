package rcache

import (
	"fmt"
	"strings"
	"testing"

	"bitbucket.org/maksadbek/go-mon-service/conf"
)

func TestInit(t *testing.T) {
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
	err = Initialize(app)
	if err != nil {
		t.Error(err)
	}
	test := "bar"
	rc.Send("SET", "foo", test)
	rc.Send("GET", "foo")
	rc.Flush()
	rc.Receive()
	v, err := rc.Receive()
	if err != nil {
		t.Error(err)
	}
	if fmt.Sprintf("%s", v) != test {
		t.Errorf("want %s, got %s\n", v, test)
	}
	// close the connection
	defer rc.Close()
}

func TestFleetTrackers(t *testing.T) {
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
	err = Initialize(app)
	if err != nil {
		t.Error(err)
	}

	// get trackers
	trackersTest, err := GetTrackers(FleetTest.FleetName, 0, 100)
	if err != nil {
		t.Error(err)
	}

	// check tracker's id
	for index, val := range FleetTest.Trackers {
		got := trackersTest[index]
		if val != got {
			t.Errorf("want %s, got %s\n", val, got)
		}
	}

	// remove tracker data from redis
	for range FleetTest.Trackers {
		rc.Do("LPOP", FleetTest.FleetName)
	}
	defer rc.Close()
}
