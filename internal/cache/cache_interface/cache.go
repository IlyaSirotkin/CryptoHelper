package cacheinterface

import "context"

type CacheHandler interface {
	Connect(context.Context) error
	Read(context.Context, string) (float64, error)
	Write(context.Context, string, float64) error
	CloseConnection() error
}
