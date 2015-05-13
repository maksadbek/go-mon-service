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

func Initialize(c conf.App) error {
	err := rcache.Initialize(c)
	if err != nil {
		return err
	}
	return err
}
func GetPositionHandler(w http.ResponseWriter, r *http.Request) {
	fleetName, user, groups := r.PostFormValue("fleet"), r.PostFormValue("user"), r.PostFormValue("groups")
	log.Log.WithFields(logrus.Fields{
		"GET Request": "/positions",
		"fleet":       fleetName,
		"user":        user,
		"groups":      groups,
	}).Info("Request")

	trackers, err := datastore.UsrTrackers(user)
	if err != nil {
		panic(err)
	}

	var fleet rcache.Fleet
	fleet.Update = make(map[string]rcache.Pos)
	if trackers.Trackers[0] == "0" {
		fleet, err = rcache.GetPositionsByFleet(fleetName, 0, 100)
		if err != nil {
			panic(err)
		}
	} else {
		pos, err := rcache.GetPositions(trackers.Trackers)
		if err != nil {
			panic(err)
		}
		fleet.Update = pos
		fleet.Id = fleetName
	}

	jpos, err := json.Marshal(fleet)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(jpos))
}
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
