package logger

import (
	"github.com/Sirupsen/logrus"
)

var Log = logrus.New()

func Init(path string) {
	Log.Formatter = new(logrus.TextFormatter)
}
