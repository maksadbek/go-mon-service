package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/logger"
	"bitbucket.org/maksadbek/go-mon-service/route"
	"github.com/Sirupsen/logrus"
)

func main() {

	Environment := os.Getenv("GOMON")

	// config init
	f, err := ioutil.ReadFile("conf.toml")
	if err != nil {
		panic(err)
	}

	c := strings.NewReader(string(f))
	if err != nil {
		panic(err)
	}

	app, err := conf.Read(c)
	if err != nil {
		panic(err)
	}

	// Logger setup
	if Environment == "production" {
		logger.Log.Formatter = new(logrus.JSONFormatter)
	} else {
		logger.Log.Formatter = new(logrus.TextFormatter)
	}

	// route setup
	route.Initialize(app)
	http.Handle("/", http.FileServer(http.Dir("static/")))
	http.HandleFunc("/positions", route.GetPositionHandler)
	http.ListenAndServe(app.SRV.Port, nil)
}
