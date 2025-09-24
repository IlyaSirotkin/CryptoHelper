package service_interface

import (
	"cryptoHelper/internal/datasource/datasource_interface"
	"cryptoHelper/internal/display/display_interface"
)

type Service interface {
	SetOutput(display_interface.Display) error
	SetInput(datasource_interface.Datasource) error
	GetData(string) (float32, error)
	SendData(string) error
	Update() error
}
