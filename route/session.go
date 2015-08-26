package route

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/datastore"
	"bitbucket.org/maksadbek/go-mon-service/logger"
	"bitbucket.org/maksadbek/go-mon-service/metrics"
)

var (
	expTokens        = metrics.NewString("expTokens")
	uidNotFoundErr   = errors.New("uid not found")
	tokenNotFoundErr = errors.New("token not found")
)

type tokenKey struct {
	ID  string
	Key string
}

// in-memory container for user tokens
type Tokens struct {
	Tokens map[string]tokenKey
	sync.RWMutex
}

var TokenList Tokens

func (t *Tokens) Put(token string, tk tokenKey) {
	if len(t.Tokens) == 0 {
		t.Tokens = make(map[string]tokenKey)
	}
	t.Lock()
	t.Tokens[token] = tk
	t.Unlock()
}

func (t *Tokens) Get(token string) (tokenKey, bool) {
	tk, ok := t.Tokens[token]
	if !ok {
		return tk, false
	}
	return tk, true
}

// FindUid can be used to check whether uid with has already got token or not
func (t *Tokens) FindUid(uid string) (string, bool) {
	for token, usr := range t.Tokens {
		if usr.ID == uid {
			return token, true
		}
	}
	return "", false
}

func (t *Tokens) Del(token string) {
	t.Lock()
	if _, ok := t.Tokens[token]; ok {
		delete(t.Tokens, token)
	}
	t.Unlock()
}

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

// LogOutHandler handles user log out request
// deletes user token from container
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
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
	token := req["token"]
	// validate for empty string
	if token == "" {
		logger.ReqWarn(r, conf.ErrReq)
		http.Error(w, "Bad Request", 400)
		return
	}

	// delete the token from tokens container
	// lock it before deleting
	TokenList.Del(token)
}

// SignUpHandler handles user sign up request
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{}
	key := make([]byte, 64)            // key for HMAC computation
	decoder := json.NewDecoder(r.Body) // json decoder
	req := make(map[string]string)     // request params
	// decode
	err := decoder.Decode(&req)
	if err != nil {
		logger.ReqWarn(r, conf.ErrReq)
		response.Message = "Body should be a JSON object"
		message, _ := json.Marshal(response)
		http.Error(w, string(message), 400)
		return
	}
	// get login and hash
	user, hash, uid := req["user"], req["hash"], req["uid"]
	// validate for empty string
	if user == "" || hash == "" || uid == "" {
		logger.ReqWarn(r, conf.ErrReq)
		response.Message = "Validation Failed - missing field"
		message, _ := json.Marshal(response)
		http.Error(w, string(message), 400)
		return
	}
	// if user credentials are bad, then send 401 status
	if !datastore.CheckUser(user, hash) {
		logger.ReqWarn(r, conf.ErrReq)
		response.Message = "Bad User Credentials"
		message, _ := json.Marshal(response)
		http.Error(w, string(message), 401)
		return
	}
	// range over tokens, and
	// if has already got token,
	// then return old token
	if token, ok := TokenList.FindUid(uid); ok {
		response.Message = "User has already logged in"
		response.Token = token
		message, _ := json.Marshal(response)
		w.Write([]byte(message))
		return
	}
	// else, generate random key
	_, err = rand.Read(key)
	if err != nil {
		logger.ReqWarn(r, err.Error())
		response.Message = "System error, please try again"
		message, _ := json.Marshal(response)
		http.Error(w, string(message), 500)
		return
	}
	// compute new token
	token := base64.StdEncoding.EncodeToString(computeHMAC(user, base64.StdEncoding.EncodeToString(key)))
	// put token into container
	TokenList.Put(token, tokenKey{ID: uid, Key: base64.StdEncoding.EncodeToString(key)})
	// write tokens into debug var
	jtokens, _ := json.MarshalIndent(TokenList.Tokens, "\t", "")
	expTokens.Set(string(jtokens))
	// setup respose
	response.Message = "OK"
	response.Token = token
	message, _ := json.Marshal(response)
	// and send computed token
	w.Write([]byte(message))
}
