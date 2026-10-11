package Redis

import (
	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func ConnectRedis(password string) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:8879",
		Password: password,
		//IdleTimeout			: -1,
		//IdleCheckFrequency	: -1,
		PoolSize: 2,
		DB:       0,
	})
}
