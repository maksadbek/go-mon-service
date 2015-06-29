package route

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/datastore"
	"bitbucket.org/maksadbek/go-mon-service/logger"
)

type tokenKey struct {
	ID  string
	Key string
}

// in-memory container for user tokens
var tokens = make(map[string]tokenKey)

// computeHMAC can be used to compute HMAC hash of given message and key
func computeHMAC(msg, key string) []byte {
	k := []byte(key)
	h := hmac.New(sha256.New, k)
	h.Write([]byte(msg))
	return h.Sum(nil)
}

// checkMAC can be used to compare HMAC hash of given msg and key with given HMAC hash
func checkMAC(msg string, expMAC []byte, key string) bool {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(expMAC, expectedMAC)
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	key := make([]byte, 64)            // key for HMAC computation
	decoder := json.NewDecoder(r.Body) // json decoder
	req := make(map[string]string)     // request params
	// decode
	err := decoder.Decode(&req)
	if err != nil {
		logger.ReqWarn(r, conf.ErrReq)
		http.Error(w, "invalid req body format", 500)
		return
	}
	// get login and hash
	user, hash, uid := req["user"], req["hash"], req["uid"]
	// validate for empty string
	if user == "" || hash == "" || uid == "" {
		logger.ReqWarn(r, conf.ErrReq)
		http.Error(w, "Bad Request", 400)
		return
	}
	// if user credentials are bad, then send 400 status
	if !datastore.CheckUser(user, hash) {
		logger.ReqWarn(r, conf.ErrReq)
		http.Error(w, "Bad User Credentials", 400)
		return
	}
	// range over tokens, and
	// if has already got token,
	// then return old token
	for oldToken, usr := range tokens {
		if usr.ID == uid {
			w.Write([]byte(base64.StdEncoding.EncodeToString([]byte(oldToken))))
			return
		}
	}
	// else, generate random key
	_, err = rand.Read(key)
	if err != nil {
		logger.ReqWarn(r, err.Error())
		http.Error(w, "system error", 500)
		return
	}
	// compute new token
	token := string(computeHMAC(user, string(key)))
	// put token into container
	tokens[token] = tokenKey{ID: uid, Key: string(key)}
	// and send computed token
	w.Write([]byte(base64.StdEncoding.EncodeToString([]byte(token))))
}
