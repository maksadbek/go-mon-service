package datastore

import (
	"fmt"
	"strings"
	"testing"

	"bitbucket.org/maksadbek/go-mon-service/conf"
)

func TestGetTrackers(t *testing.T) {
	mockConf := `[ds]
    [ds.mysql]
        dsl = "root:toor@tcp(localhost:3306)/maxtrack"
    [ds.redis]
		host = ":6379"
		fprefix = "fleet"
        tprefix = "tracker"
	[srv]
		port = ":1234"
	[log]
		path = "info.log"
	`
	c := strings.NewReader(mockConf)
	app, err := conf.Read(c)
	if err != nil {
		t.Error(err)
	}

	// mysql setup
	err = Initialize(app.DS)
	if err != nil {
		t.Error(err)
	}
	pos, err := GetTrackers("202")
	fmt.Println(pos)
}
