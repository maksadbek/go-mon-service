package route

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/datastore"
	"bitbucket.org/maksadbek/go-mon-service/logger"
)

func computeHMAC(msg, key string) []byte {
	k := []byte(key)
	h := hmac.New(sha256.New, k)
	h.Write([]byte(msg))
	return h.Sum(nil)
}

func checkMAC(msg string, expMAC []byte, key string) bool {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(expMAC, expectedMAC)
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	req := make(map[string]string)
	err := decoder.Decode(&req)
	if err != nil {
		logger.ReqWarn(r, conf.ErrReq)
		http.Error(w, "invalid req body format", 500)
		return
	}
	user, hash := req["user"], req["hash"]
	if user == "" || hash == "" {
		logger.ReqWarn(r, conf.ErrReq)
		http.Error(w, "Bad Request", 400)
		return
	}

	if !datastore.CheckUser(user, hash) {
		logger.ReqWarn(r, conf.ErrReq)
		http.Error(w, "Bad Request", 400)
		return
	}

	token := computeHMAC(user, config.Auth.MACKey)
	w.Write([]byte(base64.StdEncoding.EncodeToString(token)))
}
