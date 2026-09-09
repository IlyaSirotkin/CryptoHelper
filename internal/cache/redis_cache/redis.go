package rediscache

type RedisHadler struct {
}

func NewRedisHandler() *RedisHadler {
	return &RedisHadler{}
}

func (rds *RedisHadler) Connect() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // use default DB
		Protocol: 2,
	})
}
