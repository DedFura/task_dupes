package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func SetupLog(filePath string) (*logrus.Logger, error) {
	var log = logrus.New()
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	log.Out = file
	log.Formatter = &logrus.JSONFormatter{}
	log.SetLevel(logrus.InfoLevel)

	return log, nil
}
