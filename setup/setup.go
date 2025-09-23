package setup

import (
	logger "cryptoHelper/pkg/applogger"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func SetENVreading(envFilePath string) error {
	err := godotenv.Load(envFilePath)
	if err != nil {
		return fmt.Errorf("error with godotenv.Load(): %w", err)
	} else {
		return nil
	}
}

func SetLogger() error {
	err := logger.Get().SetOutputFile(os.Getenv("LOG_FILE_NAME"))
	if err != nil {
		return fmt.Errorf("error with logger's output file opening: %w", err)
	} else {
		return nil
	}
}
