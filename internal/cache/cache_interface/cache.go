package cacheinterface

type CacheHandler interface {
	Connect()
	Read()
	Write()
	Close()
}
