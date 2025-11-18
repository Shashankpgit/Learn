package utils

import (
	"app/src/constants"
	"os"

	"github.com/sirupsen/logrus"
)

type CustomFormatter struct {
	logrus.TextFormatter
}

// NewLogger creates and configures a new logger instance
// This is a constructor function for dependency injection
func NewLogger() *logrus.Logger {
	log := logrus.New()

	// Set logger to use the custom text formatter
	log.SetFormatter(&CustomFormatter{
		TextFormatter: logrus.TextFormatter{
			TimestampFormat: constants.LoggerTimestampFormat,
			FullTimestamp:   true,
			ForceColors:     true,
		},
	})

	log.SetOutput(os.Stdout)
	return log
}
