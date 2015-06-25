package route

import (
	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/datastore"
	"bitbucket.org/maksadbek/go-mon-service/logger"
	"bitbucket.org/maksadbek/go-mon-service/rcache"
	"github.com/garyburd/redigo/redis"
)

var config conf.App

func Initialize(c conf.App) error {
	config = c
	err := rcache.Initialize(config)
	if err != nil {
		return err
	}
	return err
}

// GetTrackers can be used to get list of trackers
// if user does not exist in cache then in caches from mysql
func GetTrackers(name string) (trackers rcache.Usr, err error) {
	trackers, err = rcache.UsrTrackers(name)
	logger.FuncLog("route.GetTracker", "GetTracker", nil, nil)
	if err == nil || err != redis.ErrNil {
		return
	}
	// if redis result is nil
	trackers, err = datastore.UsrTrackers(name)
	if err != nil {
		return
	}
	err = rcache.SetUsrTrackers(trackers)
	if err != nil {
		logger.FuncLog("route.GetTrackers", "GetTrackers", nil, err)
		return
	}

	return
}
