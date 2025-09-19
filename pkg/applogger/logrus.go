package applogger

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

type applogger struct{}

var (
	instance *applogger
	once     sync.Once
)

func GetLogger() *applogger {
	once.Do(func() {
		log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
		instance = &applogger{}
	})
	return instance
}

func (l *applogger) Info(msg string) {
	log.Info(msg)
}

func (l *applogger) Debug(msg string) {
	log.Debug(msg)
}

func (l *applogger) Warning(msg string) {
	log.Warn(msg)
}

func (l *applogger) Error(msg string) {
	log.Error(msg)
}

func (l *applogger) Fatal(msg string) {
	log.Fatal(msg)
}

func (l *applogger) SetOutputFile(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	log.SetOutput(file)
	return nil
}
