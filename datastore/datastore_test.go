package datastore

import (
	"strings"
	"testing"

	"bitbucket.org/maksadbek/go-mon-service/conf"
)

func TestGetTrackers(t *testing.T) {
	c := strings.NewReader(mockConf)
	app, err := conf.Read(c)
	if err != nil {
		t.Error(err)
	}

	// mysql setup
	err = Initialize(app)
	if err != nil {
		t.Error(err)
	}
	_, err = GetTrackers("202")
}

func TestUsrTrackersPartialCars(t *testing.T) {
	usr, err := UsrTrackers("Kamilka")
	if err != nil {
		t.Error(err)
	}
	if usr.Trackers[0] == "0" {
		t.Errorf("want %s, got %s", "0", usr.Trackers[0])
	}
}

func TestUsrTrackersAllCars(t *testing.T) {
	usr, err := UsrTrackers("newmax")
	if err != nil {
		t.Error(err)
	}
	if usr.Trackers[0] != "0" {
		t.Errorf("want %s, got %s", "0", usr.Trackers[0])
	}
}

func TestCacheFleetTrackers(t *testing.T) {
	_, err := CacheFleetTrackers()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkGetTrackers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetTrackers("")
	}
}

func BenchmarkCacheFleetTrackers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = CacheFleetTrackers()
	}
}

func TestGetLitrage(t *testing.T) {
	_, err := GetLitrage(104953, 40)
	if err != nil {
		t.Error(err)
	}
}

func TestCheckUser(t *testing.T) {
	res := CheckUser(UserTest.Username, UserTest.Hash)
	if !res {
		t.Errorf("want %t, got %t", true, res)
	}
}
