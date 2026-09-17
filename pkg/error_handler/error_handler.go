package error_handler

import (
	logger "cryptoHelper/pkg/applogger"
	"os"
)

func ErrorCatch(err error, msg string) {
	if err != nil {
		logger.Get().Error(msg + err.Error())
		os.Exit(1)
	}
}
