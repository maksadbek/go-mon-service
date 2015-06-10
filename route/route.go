package route

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/datastore"
	"bitbucket.org/maksadbek/go-mon-service/logger"
	"bitbucket.org/maksadbek/go-mon-service/rcache"
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
	decoder := json.NewDecoder(r.Body)
	req := make(map[string]string)
	err := decoder.Decode(&req)
	if err != nil {
		logger.ReqWarn(r, conf.ErrReq)
		http.Error(w, "invalid req body format", 500)
		return
	}
	fmt.Printf("%+v\n", req)
	//fleetName, user, groups := r.PostFormValue("selectedFleetJs"), r.PostFormValue("user"), r.PostFormValue("groups")
	fleetName, user, groups := req["selectedFleetJs"], req["user"], req["groups"]
	if fleetName == "" || user == "" || groups == "" {
		logger.ReqWarn(r, conf.ErrReq)
		http.Error(w, config.ErrorMsg["NotExistInCache"].Msg, 404)
		return
	}
	logger.ReqWarn(r, conf.ErrReq)
	trackers, err := GetTrackers(user)
	if err != nil {
		logger.ReqWarn(r, conf.ErrReq, err)
		http.Error(w, err.Error(), 404)
		return
	}
	var fleet rcache.Fleet
	fleet.Update = make(map[string]rcache.Pos)
	if trackers.Trackers[0] == "0" {
		fleet, err = rcache.GetPositionsByFleet(fleetName, 0, 100)
		if err != nil {
			logger.ReqWarn(r, conf.ErrReq, err)
			http.Error(w, err.Error(), 404)
			return
		}
	} else {
		pos, err := rcache.GetPositions(trackers.Trackers)
		if err != nil {
			logger.ReqWarn(r, conf.ErrReq, err)
			http.Error(w, err.Error(), 404)
			return
		}
		fleet.Update = pos
		fleet.Id = fleetName
	}

	jpos, err := json.Marshal(fleet)
	if err != nil {
		logger.ReqWarn(r, conf.ErrReq, err)
		http.Error(w, err.Error(), 404)
		return
	}
	w.Write([]byte(jpos))
	return
}

func GetTrackers(name string) (rcache.Usr, error) {
	trackers, err := rcache.UsrTrackers(name)
	logger.FuncLog("route.GetTracker", "GetTracker", nil, nil)
	if err != nil {
		if err.Error() == config.ErrorMsg["NotExistInCache"].Msg {
			trackersDS, err := datastore.UsrTrackers(name)
			if err != nil {
				return trackers, err
			}
			err = rcache.SetUsrTrackers(trackersDS)
			if err != nil {
				logger.FuncLog("route.GetTrackers", "GetTrackers", nil, err)
				return trackersDS, err
			}
			return trackersDS, nil
		} else {
			return trackers, err
		}
	}
	return trackers, nil
}
