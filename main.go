package main

import (
	_ "expvar"
	"flag"
	"net"
	"net/http"
	"os"
	"runtime"
	"syscall"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/datastore"
	log "bitbucket.org/maksadbek/go-mon-service/logger"
	"bitbucket.org/maksadbek/go-mon-service/metrics"
	"bitbucket.org/maksadbek/go-mon-service/rcache"
	"bitbucket.org/maksadbek/go-mon-service/route"
	"bitbucket.org/maksadbek/go-mon-service/utils"
	"github.com/Sirupsen/logrus"
	"github.com/rs/cors"
	"github.com/sevlyar/go-daemon"
)

var (
	confPath = flag.String(
		"conf",
		"conf.toml",
		`conf file for daemon`)
	logPath = flag.String(
		"log",
		"log",
		`log file path for daemon`)
	signal = flag.String("s", "",
		`send signal to the daemon
		quit — graceful shutdown
		stop — fast shutdown
		reload — reloading the configuration file`)
	logLevel = flag.String("v", "error", "log level: debug, info, warn, error")
	done     = make(chan struct{})
	app      conf.App
)

type Server struct {
	Listener net.Listener
}

func (srv *Server) Close() {
	srv.Listener.Close()
}

var server Server

func main() {
	daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)

	runtime.GOMAXPROCS(4)
	// parse flags
	flag.Parse()
	// set maxprocs to 4
	cntxt := &daemon.Context{
		PidFileName: "pid",
		PidFilePerm: 0644,
		LogFileName: *logPath,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}
	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Log.Fatal("Unable send signal to the daemon:", err)
		}
		daemon.SendCommands(d)
		return
	}
	d, err := cntxt.Reborn()
	if err != nil {
		log.Log.Fatal(err)
	}
	log.Log.Info("daemon started")
	if d != nil {
		return
	}
	defer cntxt.Release()
	log.Log.Info("daemon started")

	// initialize log level
	log.Init(*logLevel)
	log.Log.Formatter = new(logrus.TextFormatter)

	app, err := conf.Init(*confPath)
	if err != nil {
		log.Log.Error(err)
	}
	go worker(app)
	err = daemon.ServeSignals()
	if err != nil {
		log.Log.Fatal(err)
	}
}

// the main worker func that turn on web server
func worker(app conf.App) {
	route.Initialize(app)
	datastore.Initialize(app)
	rcache.Initialize(app)
	go utils.CacheData(app, done)
	go utils.CacheGroups(app, done)
	// setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"content-type", "x-requested-with"},
		AllowedMethods:   []string{"post"},
	})
	var err error
	server.Listener, err = net.Listen("tcp", app.SRV.Port)
	if err != nil {
		log.Log.Error(err)
	}
	handler := c.Handler(webHandlers())
	serve := &http.Server{Handler: route.GzipHandler(handler)}
	serve.Serve(server.Listener)
	<-done
}

func webHandlers() http.Handler {
	web := http.NewServeMux()
	web.Handle("/", http.FileServer(http.Dir("static/")))
	web.HandleFunc("/positions", route.GetPositionHandler)
	web.HandleFunc("/signup", route.SignupHandler)
	web.HandleFunc("/logout", route.LogoutHandler)
	web.HandleFunc("/debug/vars/", metrics.MetricsHandler)
	//metrics.Publish("cmdline", metrics.Func(metrics.Cmdline))
	metrics.Publish("memstats", metrics.Func(metrics.Memstats))
	metrics.Publish("goroutines", metrics.Func(metrics.Goroutines))
	return web
}

func termHandler(sig os.Signal) error {
	log.Log.Info("terminating...")
	server.Close()
	done <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
	log.Log.Info("reloaded")
	return nil
}
