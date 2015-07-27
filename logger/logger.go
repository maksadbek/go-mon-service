package logger

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

var Log = logrus.New()

func Init(level string) {
	Log.Formatter = new(logrus.TextFormatter)
	switch {
	case level == "debug":
		Log.Level = logrus.DebugLevel
	case level == "info":
		Log.Level = logrus.InfoLevel
	case level == "warn":
		Log.Level = logrus.WarnLevel
	case level == "fatal":
		Log.Level = logrus.FatalLevel
	case level == "panic":
		Log.Level = logrus.PanicLevel
	}
}

func ReqWarn(req *http.Request, msg string, err ...error) {
	fleetName, user, groups := req.PostFormValue("fleet"), req.PostFormValue("user"), req.PostFormValue("groups")
	// prepare map
	m := make(map[string]interface{})
	m["method"] = req.URL.Path
	m["fleet"] = fleetName
	m["user"] = user
	m["groups"] = groups
	for _, e := range err {
		if e != nil {
			m["error"] = e.Error()
		}
	}
	m["http status"] = 404
	Log.WithFields(logrus.Fields{"": m}).Warn(msg)
}

func FuncLog(fn, msg string, msgs map[string]interface{}, err error) {
	/*
		m := make(map[string]interface{})
		if msgs != nil {
			m = msgs
		}
		m["package"] = fn

		if err != nil {
			m["error"] = err.Error()
			Log.WithFields(logrus.Fields{"": fmt.Sprintf("%+v", m)}).Warn(msg)
		} else {
			Log.WithFields(logrus.Fields{"": fmt.Sprintf("%+v", m)}).Debug(msg)
		}
	*/
}
