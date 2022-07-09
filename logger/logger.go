package logger

import "github.com/sirupsen/logrus"

var Log *logrus.Logger

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	Log = logrus.New()
}
