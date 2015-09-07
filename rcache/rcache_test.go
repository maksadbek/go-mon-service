package rcache

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"bitbucket.org/maksadbek/go-mon-service/conf"
)

func TestInit(t *testing.T) {
	// close the connection
	// defer rc.Close()
	r := strings.NewReader(mockConf)
	app, err := conf.Read(r)
	err = Initialize(app)
	if err != nil {
		t.Error(err)
	}
}

func TestFleetTrackers(t *testing.T) {
	r := strings.NewReader(mockConf)
	app, err := conf.Read(r)
	err = Initialize(app)
	if err != nil {
		t.Error(err)
	}

	// get trackers
	trackersTest, err := GetTrackers(
		FleetTest.FleetName,
		0,
		100,
	)
	if err != nil {
		t.Error(err)
	}

	// check tracker's id
	for _, track := range trackersTest {
		found := false
		for _, val := range FleetTest.Trackers {
			if val == track {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%s not found", track)
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
	t.Logf("%+v\n", flt)
	/*
		for _, x := range FleetTest.Trackers {
			if flt.Update[x].Id != testFleet.Update[x].Id {
				t.Errorf("want %+v, got %+v", testFleet.Update[x], flt.Update[x])
			}
		}
	*/
}

func TestGetPostions(t *testing.T) {
	pos, err := GetPositions(FleetTest.Trackers)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v\n", pos)
	/*
		for _, tracker := range pos {
			idStr := strconv.Itoa(tracker.Id)
			if idStr != FleetTest.Trackers[0] &&
				idStr != FleetTest.Trackers[1] &&
				idStr != FleetTest.Trackers[2] {
				t.Errorf(
					"want %s or %s or %s, got %s\n",
					FleetTest.Trackers[0],
					FleetTest.Trackers[1],
					FleetTest.Trackers[2],
					idStr,
				)
			}
		}
	*/
}

func TestUsrTrackers(t *testing.T) {
	usr, err := UsrTrackers(testUsr[0].Login)
	if err != nil {
		t.Error(err)
	}

	want := testUsr[0].Login
	if usr.Login != want {
		t.Errorf("want %s, got %s", want, usr.Login)
	}
}

func TestSetUsrTrackers(t *testing.T) {
	err := SetUsrTrackers(testUsr[1])
	if err != nil {
		t.Error(err)
	}
	rc := pool.Get()
	defer rc.Close()
	userb, err := rc.Do(
		"GET",
		config.DS.Redis.UPrefix+":"+testUsr[1].Login,
	)
	if fmt.Sprintf("%v", userb) == "<nil>" {
		t.Error("got nil")
	}
	usr := Usr{}
	err = json.Unmarshal([]byte(fmt.Sprintf("%s", userb)), &usr)
	if err != nil {
		t.Error(err)
	}

	want := testUsr[1].Login
	if usr.Login != want {
		t.Errorf("got %s, want %s", usr.Login, want)
	}
}

// test for non existing tracker data
func TestGetPositions_NonExisting(t *testing.T) {
	v, err := GetPositions([]string{"10", "3010"})
	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v\n", v)
	/*
		if v["10"].Id != 0 {
			t.Errorf("got %v, want nothing", v["10"])
		}
	*/
}

func TestMarshalJsonFleet(t *testing.T) {
	jfleet, err := json.MarshalIndent(testFleet, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jfleet))
}
