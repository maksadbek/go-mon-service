package route

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/rcache"
)

func Initialize(c conf.App) error {
	err := rcache.Initialize(c)
	if err != nil {
		return err
	}
	return err
}
func GetPositionHandler(w http.ResponseWriter, r *http.Request) {
	pos, err := rcache.GetPositions("fleet_202", 0, 100)
	if err != nil {
		panic(err)
	}
	jpos, err := json.Marshal(pos)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(jpos))
}
