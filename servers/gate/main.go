package main

import (
	logrus2 "github.com/go-micro/plugins/v4/logger/logrus"
	"github.com/sirupsen/logrus"
)

func main() {
	l := logrus.New() // *logrus.Logger
	logger.DefaultLogger = logger.NewLogger(logrus2.WithLogger(l))
	logger.Infof("testing: %s", "Infof")
}
