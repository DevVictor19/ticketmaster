package server

import "github.com/redis/go-redis/v9"

func getRedisClient(cfg *Redis) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return rdb
}
