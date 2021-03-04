package cache

import (
	"github.com/fatihsezgin/candlecloud-backend/internal/config"
	"github.com/go-redis/redis/v8"
)

var client *redis.Client

func InitRedis(cfg *config.RedisConfiguration) {
	client = redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})
}

func GetClient() *redis.Client {
	return client
}
