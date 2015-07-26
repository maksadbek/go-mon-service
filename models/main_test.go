package models

import (
	"os"
	"strings"
	"testing"

	"github.com/Maksadbek/wherepo/conf"
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
