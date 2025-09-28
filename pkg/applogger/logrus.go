package applogger

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type applogger struct{}

var (
	instance *applogger
	once     sync.Once
)

func Get() *applogger {
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
		return fmt.Errorf("error with log's file folder: %w", err1)
	}
	filePathName := dirName + "/" + fileName
	logrus.SetOutput(&lumberjack.Logger{
		Filename:   filePathName,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	})
	return nil
}
