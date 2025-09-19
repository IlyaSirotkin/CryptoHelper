package applogger

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type applogger struct{}

var (
	instance *applogger
	once     sync.Once
)

func GetLogger() *applogger {
	once.Do(func() {
		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
		logrus.SetLevel(logrus.DebugLevel)
		instance = &applogger{}
	})
	return instance
}

func (l *applogger) Info(msg string) {
	logrus.Info(msg)
}

func (l *applogger) Debug(msg string) {
	logrus.Debug(msg)
}

func (l *applogger) Warning(msg string) {
	logrus.Warn(msg)
}

func (l *applogger) Error(msg string) {
	logrus.Error(msg)
}

func (l *applogger) Fatal(msg string) {
	logrus.Fatal(msg)
}

func (l *applogger) SetOutputFile(fileName string) error {
	dirName := os.Getenv("LOG_DIR_NAME")
	err1 := os.MkdirAll(dirName, 0755)
	if err1 != nil && !os.IsExist(err1) {
		return fmt.Errorf("Error with log's file folder: %w", err1)
	}
	filePathName := dirName + "/" + fileName
	file, err2 := os.OpenFile(filePathName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err2 != nil {
		return fmt.Errorf("Error with log's file opening: %w", err2)
	}

	logrus.SetOutput(file)
	return nil
}
