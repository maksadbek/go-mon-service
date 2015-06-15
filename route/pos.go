package route

import (
	//	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/logger"
	"bitbucket.org/maksadbek/go-mon-service/rcache"
)

//w.Write([]byte())
func GetPositionHandler(w http.ResponseWriter, r *http.Request) {
	// decode request values
	decoder := json.NewDecoder(r.Body)
	req := make(map[string]string)
	err := decoder.Decode(&req)
	if err != nil {
		logger.ReqWarn(r, conf.ErrReq)
		http.Error(w, "invalid req body format", 500)
		return
	}

	// retrieve data
	fleetName := req["selectedFleetJs"]
	user := req["user"]
	groups := req["groups"]
	token := "token"

	// validate for empty
	if fleetName == "" || user == "" || groups == "" || token == "" {
		logger.ReqWarn(r, conf.ErrReq)
		http.Error(w, conf.ErrInvalidReq, 400)
		return
	}

	// check token
	//expectedToken := computeHMAC(user, config.Auth.MACKey)
	// if token != base64.StdEncoding.EncodeToString(expectedToken) {
	//		http.Error(w, conf.ErrUnauthReq, 511)
	//		return
	//	}
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

	fleet.LastReq = time.Now().Unix()
	jpos, err := json.Marshal(fleet)
	if err != nil {
		logger.ReqWarn(r, conf.ErrReq, err)
		http.Error(w, err.Error(), 404)
		return
	}
	w.Write([]byte(jpos))
	return
}
