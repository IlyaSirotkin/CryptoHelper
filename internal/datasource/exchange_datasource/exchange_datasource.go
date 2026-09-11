package exchange_datasource

import (
	logger "cryptoHelper/pkg/applogger"
	"encoding/json"
	"net/http"
	"strconv"
)

type currentPriceSerialization struct {
	CoinName string `json:"symbol"`
	PriceUsd string `json:"price"`
}

type Exchange struct {
	exchageURL string
}

func NewExchange(url_env string) *Exchange {
	return &Exchange{exchageURL: url_env}
}

func (ex *Exchange) ExtractCurrentPrice(currencyName string) (float32, error) {

	responce, err := http.Get(ex.exchageURL + currencyName + "USDT")

	if err != nil {
		logger.Get().Debug("http requst finished with error ")
		return 0.0, err
	}
	defer responce.Body.Close()

	receivedData := &currentPriceSerialization{}
	err = json.NewDecoder(responce.Body).Decode(receivedData)
	if err != nil {
		logger.Get().Debug("http requst decoding has gone wrong ")
		return 0.0, err
	}

	price, err := strconv.ParseFloat(receivedData.PriceUsd, 32)
	if err != nil {
		logger.Get().Debug("string price has problem with float parsing ")
		return 0.0, err
	}

	logger.Get().Info(currencyName + " prices was successfully extracted")

	return float32(price), nil
}
