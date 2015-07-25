package datastore

import (
	"os"
	"strings"
	"testing"

	"bitbucket.org/maksadbek/go-mon-service/conf"
)

func TestMain(m *testing.M) {
	c := strings.NewReader(mockConf)
	app, err := conf.Read(c)
	if err != nil {
		panic(err)
	}

	// mysql init
	err = Initialize(app)
	if err != nil {
		panic(err)
	}
	retCode := m.Run()

	os.Exit(retCode)
}
