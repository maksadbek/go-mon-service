package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/datastore"
	log "bitbucket.org/maksadbek/go-mon-service/logger"
	"bitbucket.org/maksadbek/go-mon-service/route"
	"github.com/Sirupsen/logrus"
)

var (
	control = flag.String(
		"s",
		"start",
		`send signal to the daemon
         start: start daemon
         stop: shutdown
         restart: reload `)
	confPath = flag.String(
		"conf",
		"conf.toml",
		`configuration file for daemon`)
	daemon = flag.Bool("d", true, "do not touch it")
	sig  = make(chan os.Signal)
	stop = make(chan bool)
	res  = make(chan bool)
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

var server Server

func main() {
	flag.Parse()
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
	case "start":
		break
	default:
		fmt.Println(`  send signal to the daemon
                       stop  - stop the daemon
                       restart - restart the daemon`)
		os.Exit(1)
	}
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGHUP)
	go sigCatch()

	go stopd()
	go resd()
	if *daemon == true {
		res <- true
		stop <- true
	}
	Environment := os.Getenv("GOMON")

	// config init
	f, err := ioutil.ReadFile(*confPath)
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
    go CacheData(app)
	go worker(app)
	err = WritePid()
	if err != nil {
		panic(err)
	}
	<-stop
}

func WritePid() error {
	f, err := os.Create("pid")
	if err != nil {
		panic(err)
	}
	pid := os.Getpid()
	log.Log.Info("my pid is", pid)
	pidStr := strconv.Itoa(pid)
	_, err = f.Write([]byte(pidStr))
	if err != nil {
		return err
	}
	f.Close()
	return nil
}
func CacheData(app conf.App) {
	datastore.CacheData()
	for _ = range time.Tick(
		time.Duration(app.DS.Mysql.Interval) * time.Minute) {
		datastore.CacheData()
	}
}
func sigCatch() {
	for {
		select {
		case s := <-sig:
			if s == syscall.SIGTERM {
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
	log.Log.Info("terminating")
	os.Exit(0)
}

func resd() {
	<-res
	log.Log.Info("restarting")
	logFile, err := os.Create("logs")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	cmd := exec.Command(os.Args[0], "-d=false")
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	stop <- true
}

func worker(app conf.App) {
	var err error
	server.Listener, err = net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	serve := &http.Server{Handler: webHandlers()}
	serve.Serve(server.Listener)
}

func webHandlers() http.Handler {
	web := http.NewServeMux()
	web.Handle("/", http.FileServer(http.Dir("static/")))
	web.HandleFunc("/positions", route.GetPositionHandler)
	return web
}
