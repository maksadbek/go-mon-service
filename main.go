package main

import (
	_ "expvar"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strconv"
	"syscall"
    "time"

	"github.com/Maksadbek/wherepo/conf"
	"github.com/Maksadbek/wherepo/models"
	log "github.com/Maksadbek/wherepo/logger"
	"github.com/Maksadbek/wherepo/metrics"
	"github.com/Maksadbek/wherepo/cache"
	"github.com/Maksadbek/wherepo/route"
	"github.com/Maksadbek/wherepo/utils"
	"github.com/Sirupsen/logrus"
	"github.com/rs/cors"
)

var (
	profiling = flag.Bool(
		"p",
		false,
		`profiling file`,
	)
	control = flag.String(
		"s",
		"start",
		`send signal to the daemon
         start: start daemon
         stop: shutdown
         restart: reload 
         status: check the status`)
	confPath = flag.String(
		"conf",
		"conf.toml",
		`configuration file for daemon`)
	daemon    = flag.Bool("d", false, "do not touch it")
	daemonize = flag.Bool("f", true, "daemonize or not")
	logLevel  = flag.String("v", "error", "log level: debug, info, warn, error")
	sig       = make(chan os.Signal)
	stop      = make(chan bool)
	res       = make(chan bool)
	app       conf.App
)

type Server struct {
	Listener net.Listener
}

func (srv *Server) Close() {
	srv.Listener.Close()
}

var server Server

func main() {
	// set maxprocs to 4
	runtime.GOMAXPROCS(4)
	// parse flags
	flag.Parse()
	// initialize log level
	log.Init(*logLevel)
	// profile
	if *profiling {
		prof, err := os.Create("profiling.pprof")
		if err != nil {
			log.Log.Error(err)
		}
		pprof.StartCPUProfile(prof)
	}

	switch *control {
	case "stop":
		pid, err := utils.ReadPid("pid")
		if err != nil {
			log.Log.Info("cannot read pid")
			os.Exit(1)
		}
		err = utils.SendTERM(pid)
		if err != nil {
			log.Log.Info("cannot send term signal to pid:", strconv.Itoa(pid))
			os.Exit(1)
		}
		os.Exit(0)

	case "restart":
		pid, err := utils.ReadPid("pid")
		if err != nil {
			log.Log.Info("cannot read pid")
			os.Exit(1)
		}
		err = utils.SendHUP(pid)
		if err != nil {
			log.Log.Info("cannot send term signal to pid:", strconv.Itoa(pid))
			os.Exit(1)
		}
		os.Exit(0)
	case "status":
		if utils.CheckPidFile("pid") {
			pid, err := utils.ReadPid("pid")
			if err != nil {
				log.Log.Error(err)
			}
			fmt.Println("daemon is running, pid is: ", strconv.Itoa(pid))
		} else {
			fmt.Println("daemon is not running")
		}
		os.Exit(0)

	case "start":
		if utils.CheckPidFile("pid") {
			fmt.Println("daemon is already running, stop it or restart")
			os.Exit(0)
		}
		break
	default:
		fmt.Println(`  
                        send : signal to the daemon
                        stop : stop the daemon
                        restart : restart the daemon
                        status : check the status
                        `)
		os.Exit(1)
	}
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	go sigCatch()
	go stopd()
	go resd()
	if *daemon == true {
		if *daemonize == true {
			res <- true
			stop <- true
		}
	}

	// config init
	app, err := conf.Init(*confPath)
	if err != nil {
		log.Log.Error(err)
	}

	// log setup, firstly get GOMON env variable
	Environment := os.Getenv("GOMON")
	if Environment == "production" {
		log.Log.Formatter = new(logrus.JSONFormatter)
	} else {
		log.Log.Formatter = new(logrus.TextFormatter)
	}

	route.Initialize(app)
	models.Initialize(app)
	cache.Initialize(app)
	go CacheData(app)
	go CacheGroups()
	go worker(app, stop)
	err = WritePid()
	if err != nil {
		log.Log.Error(err)
	}
	<-stop
}

func WritePid() error {
	f, err := os.Create("pid")
	if err != nil {
		log.Log.Error(err)
	}
	pid := os.Getpid()
	log.Log.Info("my pid is", pid)
	pidStr := strconv.Itoa(pid)
	_, err = f.Write([]byte(pidStr))
	if err != nil {
		log.Log.Error(err)
	}
	f.Close()
	return nil
}
func CacheData(app conf.App) {
	trackers, err := models.GetTrackers()
	if err != nil {
		log.Log.Error(err)
	}
	err = cache.CacheDefaults(trackers)
	if err != nil {
		log.Log.Error(err)
	}
	CacheFleetTrackers()
	for _ = range time.Tick(time.Duration(app.DS.Mysql.Interval) * time.Minute) {
		trackers, err := models.GetTrackers()
		if err != nil {
			log.Log.Error(err)
		}
		cache.CacheDefaults(trackers)
		CacheFleetTrackers()
	}
}

func CacheFleetTrackers() {
	t, err := models.CacheFleetTrackers()
	if err != nil {

		log.Log.Error(err)
	}
	err = cache.AddFleetTrackers(t)
	if err != nil {
		log.Log.Error(err)
	}
}

func CacheGroups() {
	err := models.LoadGroups()
	if err != nil {
		log.Log.Error(err)
	}
	for _ = range time.Tick(time.Duration(app.Cache.GroupInterval) * time.Minute) {
		err := models.LoadGroups()
		if err != nil {
			log.Log.Error(err)
		}
	}
}
func sigCatch() {
	for {
		select {
		case s := <-sig:
			if s == syscall.SIGTERM || s == syscall.SIGINT {
				if *profiling {
					pprof.StopCPUProfile()
				}
				log.Log.Info("got term signal")
				stop <- true
				break
			} else if s == syscall.SIGHUP {
				log.Log.Info("got hup signal")
				server.Close()
				res <- true
				stop <- true
				break
			}

		}
	}
}
func stopd() {
	<-stop
	if utils.CheckPidFile("pid") {
		err := os.Remove("pid")
		if err != nil {
			panic(err)
		}
	}

	log.Log.Info("terminating")
	os.Exit(0)
}

func resd() {
	<-res
	log.Log.Info("restarting")
	logFile, err := os.Create("logs")
	if err != nil {
		log.Log.Error(err)
	}
	defer logFile.Close()
	cmd := exec.Command(os.Args[0], "-d=false")
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	err = cmd.Start()
	if err != nil {
		log.Log.Error(err)
	}
	// remove pid file
	if utils.CheckPidFile("pid") {
		err := os.Remove("pid")
		if err != nil {
			log.Log.Error(err)
		}
	}
}

// the main worker func that turn on web server
func worker(app conf.App, done <-chan bool) {
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
	web.Handle("/", http.FileServer(http.Dir("public/")))
	web.HandleFunc("/positions", route.GetPositionHandler)
	web.HandleFunc("/signup", route.SignupHandler)
	web.HandleFunc("/logout", route.LogoutHandler)
	web.HandleFunc("/debug/vars/", metrics.MetricsHandler)
	//metrics.Publish("cmdline", metrics.Func(metrics.Cmdline))
	metrics.Publish("memstats", metrics.Func(metrics.Memstats))
	metrics.Publish("goroutines", metrics.Func(metrics.Goroutines))
	return web
}
