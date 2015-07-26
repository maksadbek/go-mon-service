package route

import (
	//	"encoding/base64"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Maksadbek/wherepo/conf"
	"github.com/Maksadbek/wherepo/logger"
	"github.com/Maksadbek/wherepo/cache"
)

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
	token := req["token"]

	// validate for empty
	if fleetName == "" || user == "" || groups == "" || token == "" {
		logger.ReqWarn(r, conf.ErrReq)
		http.Error(w, conf.ErrInvalidReq, 400)
		return
	}

	// check token
	usrTokenKey, ok := tokenList.Get(token)
	if !ok {
		http.Error(w, conf.ErrUnauthReq, 511)
		return
	}
	expectedToken := computeHMAC(user, usrTokenKey.Key)
	if token != base64.StdEncoding.EncodeToString(expectedToken) {
		http.Error(w, conf.ErrUnauthReq, 511)
		return
	}
	logger.ReqWarn(r, conf.ErrReq)
	trackers, err := GetTrackers(user)
	if err != nil {
		logger.ReqWarn(r, conf.ErrReq, err)
		http.Error(w, err.Error(), 404)
		return
	}
	var fleet cache.Fleet
	fleet.Update = make(map[string][]cache.Pos)
	if trackers.Trackers[0] == "0" {
		fleet, err = cache.GetPositionsByFleet(fleetName, 0, 100)
		if err != nil {
			logger.ReqWarn(r, conf.ErrReq, err)
			http.Error(w, err.Error(), 404)
			return
		}
	} else {
		pos, err := cache.GetPositions(trackers.Trackers)
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
