package logging

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetReportCaller(false)
	logger.Formatter = &logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
	}

	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)
}

func GetLogger() *logrus.Logger {
	return logger
}
