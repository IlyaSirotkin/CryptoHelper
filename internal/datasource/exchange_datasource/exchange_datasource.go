package exchange_datasource

import (
	logger "cryptoHelper/pkg/applogger"
	"encoding/json"
	"fmt"
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

func NewExchange() *Exchange {
	return &Exchange{}
}

func (ex *Exchange) ExtractCurrentPrice(currencyName string) (float32, error) {

	responce, err := http.Get(os.Getenv("BINANCE_COIN_API") + currencyName + "USDT")

	if err != nil {
		logger.Get().Error("http requst finished with error: " + fmt.Sprint(err))
		return 0.0, err
	}
	defer responce.Body.Close()

	receivedData := &currentPriceSerialization{}
	err = json.NewDecoder(responce.Body).Decode(receivedData)
	if err != nil {
		logger.Get().Error("http requst decoding has gone wrong: " + fmt.Sprint(err))
		return 0.0, err
	}

	price, err := strconv.ParseFloat(receivedData.PriceUsd, 32)
	if err != nil {
		logger.Get().Error("string price has problem with float parsing: " + fmt.Sprint(err))
		return 0.0, err
	}

	logger.Get().Info(currencyName + " prices was successfully extracted")

	return float32(price), nil
}
