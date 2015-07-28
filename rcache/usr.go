package rcache

import (
	"encoding/json"

	"bitbucket.org/maksadbek/go-mon-service/logger"
	"github.com/garyburd/redigo/redis"
)

// structure for user
type Usr struct {
	Login    string   // login of user
	Fleet    string   // user's fleet id
	Trackers []string // user's list of trackers
}

// UsrTrackers can be used to get info of user and list of its trackers
func UsrTrackers(name string) (Usr, error) {
	usr := Usr{}
	rc := pool.Get()
	defer rc.Close()
	logger.FuncLog("rcache.UsrTrackers", "", nil, nil)
	// get user data
	userb, err := redis.String(rc.Do("GET", config.DS.Redis.UPrefix+":"+name)) // prefix can be set from conf
	if err != nil {
		logger.FuncLog("rcache.UsrTrackers", "", nil, err)
		return usr, err
	}
	err = json.Unmarshal([]byte(userb), &usr)
	if err != nil {
		logger.FuncLog("rcache.UsrTrackers", "", nil, err)
		return usr, err
	}
	return usr, nil
}

// SetUsrTrackers can be used to save user info in redis
func SetUsrTrackers(usr Usr) error {
	rc := pool.Get()
	defer rc.Close()
	jusr, err := json.Marshal(usr)
	if err != nil {
		logger.FuncLog("rcache.SetUsrTrackers", "", nil, err)
		return err
	}
	rc.Do(
		"SET",
		config.DS.Redis.UPrefix+":"+usr.Login,
		string(jusr),
	)
	return nil
}
