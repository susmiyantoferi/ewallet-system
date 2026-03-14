package logger

import (
	"ewallet/infrastructure/config"

	"github.com/sirupsen/logrus"
)

func NewLogrus(config *config.LoggerConfig) *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.Level(config.Level))
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}
