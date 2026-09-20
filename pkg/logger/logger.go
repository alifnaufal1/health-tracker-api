package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(resolveLevel())
	return log
}

func resolveLevel() logrus.Level {
	switch strings.ToLower(os.Getenv("APP_ENV")) {
	case "production":
		return logrus.InfoLevel
	default: 
		return logrus.DebugLevel
	}
}