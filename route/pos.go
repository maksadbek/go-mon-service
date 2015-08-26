package route

import (
	//	"encoding/base64"

	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/logger"
	"bitbucket.org/maksadbek/go-mon-service/rcache"
	"github.com/didip/tollbooth"
)

func GetPositionHandler(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Message string `json:"message"`
	}{}
	tok, ok := r.Header["X-Access-Token"]
	if !ok || len(tok) == 0 {
		response.Message = "Missing Token Key"
		message, _ := json.Marshal(response)
		http.Error(w, string(message), 400)
		return
	}
	token := tok[0]

	httpError := tollbooth.LimitByKeys(Limiter, tok)
	if httpError != nil {
		response.Message = httpError.Error()
		message, _ := json.Marshal(response)
		http.Error(w, string(message), 400)
		return
	}
	// decode request values
	decoder := json.NewDecoder(r.Body)
	req := make(map[string]string)
	err := decoder.Decode(&req)
	if err != nil {
		logger.ReqWarn(r, conf.ErrReq)
		response.Message = "invalid req body format"
		message, _ := json.Marshal(response)
		http.Error(w, string(message), 400)
		return
	}

	// retrieve data
	fleetName := req["fleetID"]
	user := req["userName"]

	// validate for empty
	if fleetName == "" || user == "" {
		logger.ReqWarn(r, conf.ErrReq)
		response.Message = "Missing Fields"
		message, _ := json.Marshal(response)
		http.Error(w, string(message), 400)
		return
	}

	// check token
	usrTokenKey, ok := TokenList.Get(token)
	if !ok {
		response.Message = conf.ErrUnauthReq
		message, _ := json.Marshal(response)
		http.Error(w, string(message), 511)
		return
	}
	expectedToken := computeHMAC(user, usrTokenKey.Key)
	if token != base64.StdEncoding.EncodeToString(expectedToken) {
		response.Message = conf.ErrUnauthReq
		message, _ := json.Marshal(response)
		http.Error(w, string(message), 511)
		return
	}
	logger.ReqWarn(r, conf.ErrReq)
	trackers, err := GetTrackers(user)
	if err != nil {
		logger.ReqWarn(r, conf.ErrReq, err)
		response.Message = err.Error()
		message, _ := json.Marshal(response)
		http.Error(w, string(message), 404)
		return
	}
	var fleet rcache.Fleet
	fleet.Update = make(map[string][]rcache.Pos)
	if trackers.Trackers[0] == "0" {
		fleet, err = rcache.GetPositionsByFleet(fleetName, 0, 100)
		if err != nil {
			logger.ReqWarn(r, conf.ErrReq, err)
			response.Message = err.Error()
			message, _ := json.Marshal(response)
			http.Error(w, string(message), 404)
			return
		}
	} else {
		pos, err := rcache.GetPositions(trackers.Trackers)
		if err != nil {
			logger.ReqWarn(r, conf.ErrReq, err)
			response.Message = err.Error()
			message, _ := json.Marshal(response)
			http.Error(w, string(message), 404)
			return
		}
		fleet.Update = pos
		fleet.Id = fleetName
	}
	fleet.LastReq = time.Now().Unix()
	jpos, err := json.Marshal(fleet)
	if err != nil {
		logger.ReqWarn(r, conf.ErrReq, err)
		response.Message = err.Error()
		message, _ := json.Marshal(response)
		http.Error(w, string(message), 404)
		return
	}
	w.Write([]byte(jpos))
	return
}
