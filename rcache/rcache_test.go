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
	fleet := "fleet_202"
	trackers := []string{"id789", "id456", "id123"}

	// get trackers
	trackersTest, err := GetTrackers(fleet, 0, 100)
	if err != nil {
		t.Error(err)
	}

	// check tracker's id
	for index, val := range trackers {
		got := trackersTest[index]
		if val != got {
			t.Errorf("Want %s, got %s\n", val, got)
		}
	}

	// remove tracker data from redis
	for range trackers {
		rc.Do("LPOP", "fleet_202")
	}
	defer rc.Close()
}
