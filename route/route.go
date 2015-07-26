package route

import (
	"github.com/Maksadbek/wherepo/conf"
	"github.com/Maksadbek/wherepo/models"
	"github.com/Maksadbek/wherepo/logger"
	"github.com/Maksadbek/wherepo/cache"
	"github.com/garyburd/redigo/redis"
)

var config conf.App

func Initialize(c conf.App) error {
	config = c
	err := cache.Initialize(config)
	if err != nil {
		return err
	}
	return err
}

// GetTrackers can be used to get list of trackers
// if user does not exist in cache then in caches from mysql
func GetTrackers(name string) (trackers cache.Usr, err error) {
	trackers, err = cache.UsrTrackers(name)
	logger.FuncLog("route.GetTracker", "GetTracker", nil, nil)
	if err == nil || err != redis.ErrNil {
		return
	}
	// if redis result is nil
	trackers, err = models.UsrTrackers(name)
	if err != nil {
		return
	}
	err = cache.SetUsrTrackers(trackers)
	if err != nil {
		logger.FuncLog("route.GetTrackers", "GetTrackers", nil, err)
		return
	}

	return
}
