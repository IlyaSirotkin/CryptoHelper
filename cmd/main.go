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
		logger.Get().Error("New Telegram service return error " + fmt.Sprint(err))
		os.Exit(1)
	}

	err = service.SetInput(exchange_datasource.NewExchange())
	if err != nil {
		logger.Get().Error("SetInput exchange_datasource return error" + fmt.Sprint(err))
		os.Exit(1)
	}

	err = service.SetOutput(
		func() display_interface.Display {
			sender, err := telegram_display.NewBotSender("TELEGRAM_BOT_TOKEN")
			if err != nil {
				logger.Get().Error("Telegram display NewBotSender return error" + fmt.Sprint(err))
				os.Exit(1)
			}
			return sender
		}(),
	)
	if err != nil {
		logger.Get().Error("SetOutput return error" + fmt.Sprint(err))
		os.Exit(1)
	}

	err = service.Update()
	if err != nil {
		logger.Get().Error("Update return error" + fmt.Sprint(err))
		os.Exit(1)
	}

}
