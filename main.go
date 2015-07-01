package main

import (
	"compress/gzip"
	_ "expvar"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/datastore"
	log "bitbucket.org/maksadbek/go-mon-service/logger"
	"bitbucket.org/maksadbek/go-mon-service/metrics"
	"bitbucket.org/maksadbek/go-mon-service/rcache"
	"bitbucket.org/maksadbek/go-mon-service/route"
	"github.com/Sirupsen/logrus"
	"github.com/rs/cors"
)

var (
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
	daemon     = flag.Bool("d", false, "do not touch it")
	daemonize  = flag.Bool("f", true, "daemonize or not")
	logLevel   = flag.String("v", "error", "log level: debug, info, warn, error")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	sig        = make(chan os.Signal)
	stop       = make(chan bool)
	res        = make(chan bool)
)

type Server struct {
	Listener net.Listener
}

func (srv *Server) Close() {
	srv.Listener.Close()
}

func sendTERM(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}
	return nil
}

func sendHUP(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = process.Signal(syscall.SIGHUP)
	if err != nil {
		return err
	}
	return nil
}

func readPid(fileName string) (int, error) {
	var pid int
	p, err := ioutil.ReadFile(fileName)
	if err != nil {
		return pid, err
	}

	pid, err = strconv.Atoi(string(p))
	if err != nil {
		return pid, err
	}
	return pid, nil
}

func checkPidFile(fileName string) bool {
	_, err := os.Stat(fileName)
	if err != nil {
		return false
	}
	return true
}

var server Server

func main() {
	runtime.GOMAXPROCS(4)
	flag.Parse()
	log.Init(*logLevel)

	switch *control {
	case "stop":
		pid, err := readPid("pid")
		if err != nil {
			log.Log.Info("cannot read pid")
			os.Exit(1)
		}
		err = sendTERM(pid)
		if err != nil {
			log.Log.Info("cannot send term signal to pid:", strconv.Itoa(pid))
			os.Exit(1)
		}
		os.Exit(0)

	case "restart":
		pid, err := readPid("pid")
		if err != nil {
			log.Log.Info("cannot read pid")
			os.Exit(1)
		}
		err = sendHUP(pid)
		if err != nil {
			log.Log.Info("cannot send term signal to pid:", strconv.Itoa(pid))
			os.Exit(1)
		}
		os.Exit(0)
	case "status":
		if checkPidFile("pid") {
			pid, err := readPid("pid")
			if err != nil {
				log.Log.Error(err)
			}
			fmt.Println("daemon is running, pid is: ", strconv.Itoa(pid))
		} else {
			fmt.Println("daemon is not running")
		}
		os.Exit(0)

	case "start":
		if checkPidFile("pid") {
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
	Environment := os.Getenv("GOMON")

	// config init
	f, err := ioutil.ReadFile(*confPath)
	if err != nil {
		log.Log.Error(err)
	}

	c := strings.NewReader(string(f))
	if err != nil {
		log.Log.Error(err)
	}

	app, err := conf.Read(c)
	if err != nil {
		log.Log.Error(err)
	}

	// log setup
	if Environment == "production" {
		log.Log.Formatter = new(logrus.JSONFormatter)
	} else {
		log.Log.Formatter = new(logrus.TextFormatter)
	}

	route.Initialize(app)
	datastore.Initialize(app)
	rcache.Initialize(app)
	go CacheData(app)
	go worker(app)
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
	trackers, err := datastore.GetTrackers()
	if err != nil {
		log.Log.Error(err)
	}
	err = rcache.CacheDefaults(trackers)
	if err != nil {
		log.Log.Error(err)
	}
	CacheFleetTrackers()
	for _ = range time.Tick(time.Duration(app.DS.Mysql.Interval) * time.Minute) {
		trackers, err := datastore.GetTrackers()
		if err != nil {
			log.Log.Error(err)
		}
		rcache.CacheDefaults(trackers)
		CacheFleetTrackers()
	}
}

func CacheFleetTrackers() {
	t, err := datastore.CacheFleetTrackers()
	if err != nil {

		log.Log.Error(err)
	}
	err = rcache.AddFleetTrackers(t)
	if err != nil {
		log.Log.Error(err)
	}
}
func sigCatch() {
	for {
		select {
		case s := <-sig:
			if s == syscall.SIGTERM || s == syscall.SIGINT {
				log.Log.Info("got term signal")
				stop <- true
				break
			} else if s == syscall.SIGHUP {
				log.Log.Info("got hup signal")
				server.Close()
				res <- true
				break
			}

		}
	}
}
func stopd() {
	<-stop
	if checkPidFile("pid") {
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
	stop <- true
	// remove pid file
	if checkPidFile("pid") {
		err := os.Remove("pid")
		if err != nil {
			log.Log.Error(err)
		}
	}
}

func worker(app conf.App) {
	var err error
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"content-type", "x-requested-with"},
		AllowedMethods:   []string{"post"},
	})
	server.Listener, err = net.Listen("tcp", app.SRV.Port)
	if err != nil {
		log.Log.Error(err)
	}
	handler := c.Handler(webHandlers())
	serve := &http.Server{Handler: GzipHandler(handler)}
	serve.Serve(server.Listener)
}

func webHandlers() http.Handler {
	web := http.NewServeMux()
	web.Handle("/", http.FileServer(http.Dir("static/")))
	web.HandleFunc("/positions", route.GetPositionHandler)
	web.HandleFunc("/signup", route.SignupHandler)
	web.HandleFunc("/logout", route.LogoutHandler)
	web.HandleFunc("/debug/vars/", metrics.MetricsHandler)
	//metrics.Publish("cmdline", metrics.Func(metrics.Cmdline))
	//metrics.Publish("memstats", metrics.Func(metrics.Memstats))
	metrics.Publish("goroutines", metrics.Func(metrics.Goroutines))
	return web
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
	sniffDone bool
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.sniffDone {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", http.DetectContentType(b))
		}
		w.sniffDone = true
	}
	return w.Writer.Write(b)
}

// Wrap a http.Handler to support transparent gzip encoding.
func GzipHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Accept-Encoding")
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		h.ServeHTTP(&gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	})
}
