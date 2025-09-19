package main

import (
	"cryptoHelper/pkg/applogger"
	"cryptoHelper/pkg/logger"
	"fmt"
	//tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	fmt.Print("fdffd")
	/*bot, err := tgBotAPI.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
	}*/
	var applog logger.Logger = applogger.GetLogger()
	applog.SetOutputFile("MyfileName")
	applog.Info("dffdfdf")
}
