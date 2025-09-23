package datasource_handler

import (
	"cryptoHelper/internal/datasource/datasource_interface"
	logger "cryptoHelper/pkg/applogger"
	"errors"
)

type handler struct {
	datasource datasource_interface.Datasource
}

func NewHandler(ds datasource_interface.Datasource) *handler {
	if ds != nil {
		logger.Get().Debug("Datasource_handler was created successfully")
		return &handler{datasource: ds}
	} else {
		logger.Get().Error("Nil value has been given to the NewHandler")
		return &handler{}
	}
}

func (h handler) GetData(currencyName string) (float32, error) {
	if h.datasource != nil {
		logger.Get().Debug("Datasource_handler called ExtractCurrentPrice() successfully")
		return h.datasource.ExtractCurrentPrice(currencyName)
	} else {
		logger.Get().Error("Datasource_interface is nil, operation can not be completed")
		return 0.0, errors.New("datasource_interface is nil, operation can not be completed")
	}
}
