package models

import (
	"strings"
	"testing"

	"github.com/Maksadbek/wherepo/conf"
	"github.com/Maksadbek/wherepo/cache"
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
	_, err = GetTrackers()
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
		_, _ = GetTrackers()
	}
}

func BenchmarkCacheFleetTrackers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = CacheFleetTrackers()
	}
}

func TestCheckUser(t *testing.T) {
	res := CheckUser(UserTest.Username, UserTest.Hash)
	if !res {
		t.Errorf("want %t, got %t", true, res)
	}
}

func TestLoadGroups(t *testing.T) {
	err := LoadGroups()
	if err != nil {
		t.Error(err)
	}
	if _, err := cache.Grouplist.Get("202"); err != nil {
		t.Error(err)
	}
}
