package util

import (
	"github.com/Waitfantasy/tmq/config"
	"github.com/go-redis/redis"
)

func NewRedisClient(c *config.Config)  *redis.Client{
	return redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})
}