package main

import (
	"cryptoHelper/internal/datasource/exchange_datasource"
	"cryptoHelper/internal/display/display_interface"
	"cryptoHelper/internal/display/telegram_display"
	"cryptoHelper/internal/service/service_interface"
	"cryptoHelper/internal/service/telegram_service"
	logger "cryptoHelper/pkg/applogger"
	"cryptoHelper/pkg/error_handler"
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
	error_handler.ErrorCatch(err, "Telegram service returned error: ")

	err = service.SetInput(exchange_datasource.NewExchange())
	error_handler.ErrorCatch(err, "Service SetInput exchange returned error: ")

	err = service.SetOutput(
		func() display_interface.Display {
			sender, err := telegram_display.NewBotSender("TELEGRAM_BOT_TOKEN")
			error_handler.ErrorCatch(err, "Telegram display Bot returned error: ")
			return sender
		}(),
	)
	error_handler.ErrorCatch(err, "Service Setoutput bot returned error: ")

	err = service.Update()
	error_handler.ErrorCatch(err, "Service Update returned error: ")

}
