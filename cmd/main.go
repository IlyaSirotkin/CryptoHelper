package main

import (
	"cryptoHelper/internal/datasource/datasource_handler"
	exchange "cryptoHelper/internal/datasource/exchange_datasource"
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

	datasourceHandler := datasource_handler.NewHandler(exchange.NewExchange())
	data, _ := datasourceHandler.GetData("SOLUSDT")
	fmt.Println(data)
	/*
		if err != nil {
		}*/
}
