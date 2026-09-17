package main

import (
	rediscache "cryptoHelper/internal/cache/redis_cache"
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
	logger.Get().Info("New Telegram service succsecfuly created")

	err = service.SetInput(exchange_datasource.NewExchange(os.Getenv("BINANCE_COIN_API")))
	error_handler.ErrorCatch(err, "Service SetInput exchange returned error: ")
	logger.Get().Info("New Exchange service succsecfuly created")

	err = service.SetCache(rediscache.NewRedisHandler(os.Getenv("REDIS_ADDRESS"), os.Getenv("REDIS_PASSWORD")))
	error_handler.ErrorCatch(err, "Service SetInput exchange returned error: ")
	logger.Get().Info("New Cache service succsecfuly created")

	err = service.SetOutput(
		func() display_interface.Display {
			sender, err := telegram_display.NewBotSender(os.Getenv("TELEGRAM_BOT_TOKEN"))
			error_handler.ErrorCatch(err, "Telegram display Bot returned error: ")
			return sender
		}(),
	)
	error_handler.ErrorCatch(err, "Service Setoutput bot returned error: ")
	logger.Get().Info("New TG Display service succsecfuly created")

	err = service.Update()
	error_handler.ErrorCatch(err, "Service Update returned error: ")
	select {}
}
