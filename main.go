package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/datastore"
	log "bitbucket.org/maksadbek/go-mon-service/logger"
	"bitbucket.org/maksadbek/go-mon-service/rcache"
	"bitbucket.org/maksadbek/go-mon-service/route"
	"github.com/Sirupsen/logrus"
)

type Server struct {
	Listener net.Listener
}

func (srv *Server) Close() {
	srv.Listener.Close()
}

var server Server

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

	// log setup
	if Environment == "production" {
		log.Log.Formatter = new(logrus.JSONFormatter)
	} else {
		log.Log.Formatter = new(logrus.TextFormatter)
	}

	route.Initialize(app)
	datastore.Initialize(app.DS)
	rcache.Initialize(app)
	go CacheData(app)
	server.Listener, err = net.Listen("tcp", app.SRV.Port)
	if err != nil {
		panic(err)
	}
	serve := &http.Server{Handler: webHandlers()}
	log.Log.Info("Serving HTTP server on " + app.SRV.Port + " port")
	serve.Serve(server.Listener)
}

func CacheData(app conf.App) {
	trackers, err := datastore.GetTrackers("")
	if err != nil {
		panic(err)
	}
	rcache.CacheDefaults(trackers)
	CacheFleetTrackers()
	for _ = range time.Tick(time.Duration(app.DS.Mysql.Interval) * time.Minute) {
		trackers, err := datastore.GetTrackers("")
		if err != nil {
			panic(err)
		}
		rcache.CacheDefaults(trackers)
		CacheFleetTrackers()
	}
}

func CacheFleetTrackers() {
	t, err := datastore.CacheFleetTrackers()
	if err != nil {
		panic(err)
	}
	rcache.AddFleetTrackers(t)
}
func webHandlers() http.Handler {
	web := http.NewServeMux()
	web.Handle("/", http.FileServer(http.Dir("static/")))
	web.HandleFunc("/positions", route.GetPositionHandler)
	return web
}
