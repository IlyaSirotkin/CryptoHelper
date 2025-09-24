package main

import (
	"cryptoHelper/internal/datasource/exchange_datasource"
	"cryptoHelper/internal/display/display_interface"
	"cryptoHelper/internal/display/telegram_display"
	"cryptoHelper/internal/service/service_interface"
	"cryptoHelper/internal/service/telegram_service"
	logger "cryptoHelper/pkg/applogger"
	setup "cryptoHelper/setup"
	"fmt"
	"os"
	//  tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	err := setup.SetENVreading("config/env_file.env")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = setup.SetLogger()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		logger.Get().Debug("logger has successfully opened log file")
	}

	var service service_interface.Service
	service, err = telegram_service.NewTelegram("TELEGRAM_BOT_TOKEN")
	if err != nil {
		logger.Get().Error("Telegram service returned error: " + fmt.Sprint(err))
		os.Exit(1)
	}

	err = service.SetInput(exchange_datasource.NewExchange())
	if err != nil {
		logger.Get().Error("Service SetInput exchange returned error: " + fmt.Sprint(err))
		os.Exit(1)
	}

	service.SetOutput(
		func() display_interface.Display {
			bot, err := telegram_display.NewBot("TELEGRAM_BOT_TOKEN")
			if err != nil {
				logger.Get().Error("Telegram display Bot returned error: " + fmt.Sprint(err))
				os.Exit(1)
			}
			return bot
		}(),
	)
}
