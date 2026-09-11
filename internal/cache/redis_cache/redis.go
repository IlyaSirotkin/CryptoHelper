package rediscache

import (
	"context"
	logger "cryptoHelper/pkg/applogger"
	"errors"
	"strconv"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type RedisHadler struct {
	rdb *redis.Client
}

func NewRedisHandler(addr, password string) *RedisHadler {
	return &RedisHadler{rdb: redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
		Protocol: 2,
	})}
}

func (rds *RedisHadler) Connect(ctx context.Context) error {
	if rds.rdb == nil {
		return errors.New("RedisHandler is not initialized ")
	} else {
		err := rds.rdb.Ping(ctx).Err()
		if err != nil {
			logger.Get().Debug("Redis Connect finished with error")
			return err
		} else {
			logger.Get().Info("Redis Connect finished succsecfully ")
			return nil
		}
	}
}

func (rds *RedisHadler) Read(ctx context.Context, key string) (float64, error) {
	val, err := rds.rdb.Get(ctx, key).Result()
	if err != nil {
		logger.Get().Debug("Redis Get finished with error")
		return 0.0, err
	}
	price, err := strconv.ParseFloat(val, 32)
	if err != nil {
		logger.Get().Debug("string price has problem with float parsing ")
		return 0.0, err
	}
	logger.Get().Info("Redis REad finished succsecfully ")
	return price, nil
}

func (rds *RedisHadler) Write(ctx context.Context, key string, val float64) error {
	err := rds.rdb.Set(ctx, key, strconv.FormatFloat(val, 'f', 2, 64), 10*time.Second).Err()
	if err != nil {
		logger.Get().Debug("Redis Set finished with error")
		return err
	}
	logger.Get().Info("Redis Write finished succsecfully ")
	return nil
}

func (rds *RedisHadler) CloseConnection() error {
	err := rds.rdb.Close()
	if err != nil {
		logger.Get().Debug("Redis CLose finished with error")
		return err
	}
	return nil
}
