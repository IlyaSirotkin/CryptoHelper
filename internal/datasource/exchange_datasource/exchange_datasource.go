package exchange_datasource

import (
	logger "cryptoHelper/pkg/applogger"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

/*
type Interval int
const (

	Day Interval = iota
	Week
	Month

)
*/
type currentPriceSerialization struct {
	CoinName string `json:"symbol"`
	PriceUsd string `json:"price"`
}

type Exchange struct {
}

func (ex *Exchange) ExtractCurrentPrice(currencyName string) (float32, error) {
	responce, err := http.Get(os.Getenv("BINANCE_COIN_API") + currencyName)
	if err != nil {
		logger.Get().Fatal("http requst finished with error: " + err.Error())
	}
	defer responce.Body.Close()
	receivedData := &currentPriceSerialization{}
	if err := json.NewDecoder(responce.Body).Decode(receivedData); err != nil {
		logger.Get().Fatal("http requst decoding has gone wrong: " + err.Error())
	}

	price, err := strconv.ParseFloat(receivedData.PriceUsd, 32)
	if err != nil {
		logger.Get().Fatal("string price has problem with float parsing: " + err.Error())
	}

	logger.Get().Info(currencyName + " prices was successfully extracted")
	return float32(price), nil
}
