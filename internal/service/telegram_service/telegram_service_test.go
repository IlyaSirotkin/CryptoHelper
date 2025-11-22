package telegram_service

import (
	"cryptoHelper/internal/datasource/exchange_datasource"
	logger "cryptoHelper/pkg/applogger"
	setup "cryptoHelper/setup"
	"fmt"
	"testing"
)

func TestGetData(t *testing.T) {
	err := setup.SetENVreading("../../../config/env_file.env")
	if err != nil {
		t.Error(err)
	}
	tg, err := NewTelegram("TELEGRAM_BOT_TOKEN")
	if err != nil {
		t.Error("NewTelegram service return error " + fmt.Sprint(err))
	}

	err = tg.SetInput(exchange_datasource.NewExchange())
	if err != nil {
		logger.Get().Error("SetInput exchange_datasource return error" + fmt.Sprint(err))
	}

	currenciesToTest := []string{"BTC", "ETH", "ONDO", "XRP", "ADA", "SOL"}

	for _, el := range currenciesToTest {
		price, err := tg.GetData(el)
		if err != nil || price == 0.0 {
			t.Error("GetData return error or price is null with " + el)
		}

	}

}
