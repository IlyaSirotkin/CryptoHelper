package error_handler

import (
	logger "cryptoHelper/pkg/applogger"
	"fmt"
	"os"
)

func ErrorCatch(err error, msg string) {
	if err != nil {
		logger.Get().Error(msg + fmt.Sprint(err))
		os.Exit(1)
	}
}
