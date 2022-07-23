package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	Log = logrus.New()
	if os.Getenv("DEBUG") == "true" {
		Log.SetLevel(logrus.DebugLevel)
	} else {
		Log.SetLevel(logrus.InfoLevel)
	}
}
