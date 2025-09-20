package main

import (
	exchange "cryptoHelper/internal/datasource/exchange_datasource"
	logger "cryptoHelper/pkg/applogger"
	"fmt"
	"os"

	//  tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("config/env_file.env")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = logger.Get().SetOutputFile(os.Getenv("LOG_FILE_NAME"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		logger.Get().Debug("logger has successfully opened log file")
	}
	ex := exchange.Exchange{}
	data, _ := ex.ExtractCurrentPrice("BTCUSDT")
	fmt.Println(data)
	/*bot, err := tgBotAPI.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
	}*/
}
