package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	"bitbucket.org/maksadbek/go-mon-service/conf"
	"bitbucket.org/maksadbek/go-mon-service/route"
)

func main() {
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

	route.Initialize(app)
	http.Handle("/", http.FileServer(http.Dir("static/")))
	http.HandleFunc("/positions", route.GetPositionHandler)
	http.ListenAndServe(":8088", nil)
}
