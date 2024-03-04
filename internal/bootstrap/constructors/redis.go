package constructors

import (
	"github.com/redis/go-redis/v9"
	"github.com/serhii-marchuk/blog/internal/bootstrap/configs"
)

type RedisClient struct {
	RC *redis.Client
}

func NewRedisClient(cfg *configs.Configs) *RedisClient {
	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.GetAddr(),
		Password: cfg.Redis.PASSWORD,
		DB:       cfg.Redis.DB,
	})

	return &RedisClient{RC: rc}
}
