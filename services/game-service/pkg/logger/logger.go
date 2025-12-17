package logger

import (
	"github.com/sirupsen/logrus"
)

func New(level string) *logrus.Logger {
	log := logrus.New()
	
	// JSON formatter for structured logging
	log.SetFormatter(&logrus.JSONFormatter{})
	
	// Set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	log.SetLevel(logLevel)
	
	return log
}