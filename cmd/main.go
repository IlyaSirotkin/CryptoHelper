package main

import (
	"cryptoHelper/pkg/applogger"
	"cryptoHelper/pkg/logger"
	"fmt"
	"os"

	//  tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("config/env_file.env")
	var logger logger.Logger = applogger.GetLogger()
	err := logger.SetOutputFile(os.Getenv("LOG_FILE_NAME"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	/*bot, err := tgBotAPI.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
	}*/
}
