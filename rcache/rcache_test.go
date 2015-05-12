package rcache

import (
	"fmt"
	"strings"
	"testing"
    "strconv"

	"bitbucket.org/maksadbek/go-mon-service/conf"
)

func TestInit(t *testing.T) {
	// close the connection
	defer rc.Close()
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
}

func TestFleetTrackers(t *testing.T) {
	defer rc.Close()
	mockConf := `
[ds]
    [ds.redis]
		host = ":6379"
		fprefix = "fleet"
        tprefix = "tracker"
[srv]
	port = ":1234"
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
}

//
func TestPushToRedis(t *testing.T) {
	// push mock data into redis
	err := PushRedis(testFleet)
	if err != nil {
		t.Error(err)
	}
}

func TestGetPositionsByFleet(t *testing.T) {
	flt, err := GetPositionsByFleet(FleetTest.FleetName, 0, 100)
	if err != nil {
		t.Error(err)
	}
	for _, x := range FleetTest.Trackers {
		if flt.Update[x].Id != testFleet.Update[x].Id {
			t.Errorf("want %+v, got %+v", testFleet.Update[x], flt.Update[x])
		}
	}
}

func TestGetPostions(t *testing.T){
    pos, err := GetPositions(FleetTest.Trackers)
    if err != nil {
            t.Error(err)
    }
    for _, tracker := range pos{
            idStr := strconv.Itoa(tracker.Id)
            if(idStr != FleetTest.Trackers[0] && 
               idStr != FleetTest.Trackers[1]){
                    t.Errorf(
                            "want %s or %s, got %s\n",
                            FleetTest.Trackers[0],
                            FleetTest.Trackers[1],
                            idStr,
                    )
            }
    }
}
