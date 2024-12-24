package logutil

import (
	hzlogrus "github.com/hertz-contrib/logger/logrus"
	"github.com/sirupsen/logrus"
)

var Logger = &hzlogrus.Logger{}

func init() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	Logger = hzlogrus.NewLogger(hzlogrus.WithLogger(logger))
}
