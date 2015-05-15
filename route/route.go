package route

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/datastore"
	log "bitbucket.org/maksadbek/go-mon-service/logger"
	"bitbucket.org/maksadbek/go-mon-service/rcache"
	"github.com/Sirupsen/logrus"
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
func GetPositionHandler(w http.ResponseWriter, r *http.Request) {
	fleetName, user, groups := r.PostFormValue("fleet"), r.PostFormValue("user"), r.PostFormValue("groups")

    if fleetName == "" || user == "" || groups == "" {
        log.Log.WithFields(logrus.Fields{
            "GET Request": "/positions",
            "fleet":       fleetName,
            "user":        user,
            "groups":      groups,
            "http status": 404,
        }).Warn("Request Error")
        http.Error(w, config.ErrorMsg["NotExistInCache"].Msg, 404)
        return
    }
	log.Log.WithFields(logrus.Fields{
		"GET Request": "/positions",
		"fleet":       fleetName,
		"user":        user,
		"groups":      groups,
	}).Info("Request")

	trackers, err := datastore.UsrTrackers(user)
	if err != nil {
        log.Log.WithFields(logrus.Fields{
            "GET Request": "/positions",
            "fleet":       fleetName,
            "user":        user,
            "groups":      groups,
            "error":      err.Error(),
            "http status": 404,
        }).Warn("Request Error")
        http.Error(w, err.Error(), 404)
        return
	}

	var fleet rcache.Fleet
	fleet.Update = make(map[string]rcache.Pos)
	if trackers.Trackers[0] == "0" {
		fleet, err = rcache.GetPositionsByFleet(fleetName, 0, 100)
		if err != nil {
            log.Log.WithFields(logrus.Fields{
                "GET Request": "/positions",
                "fleet":       fleetName,
                "user":        user,
                "groups":      groups,
                "error":      err.Error(),
                "http status": 404,
            }).Warn("Request Error")
            http.Error(w, err.Error(), 404)
            return
		}
	} else {
		pos, err := rcache.GetPositions(trackers.Trackers)
		if err != nil {
            log.Log.WithFields(logrus.Fields{
                "GET Request": "/positions",
                "fleet":       fleetName,
                "user":        user,
                "groups":      groups,
                "error":      err.Error(),
                "http status": 404,
            }).Warn("Request Error")
            http.Error(w, err.Error(), 404)
            return
		}
		fleet.Update = pos
		fleet.Id = fleetName
	}

	jpos, err := json.Marshal(fleet)
	if err != nil {
        log.Log.WithFields(logrus.Fields{
            "GET Request": "/positions",
            "fleet":       fleetName,
            "user":        user,
            "groups":      groups,
            "error":      err.Error(),
            "http status": 404,
        }).Warn("Request Error")
        http.Error(w, err.Error(), 404)
        return
	}

	w.Write([]byte(jpos))
}
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
