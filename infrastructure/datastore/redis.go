package datastore

import (
	"ewallet/infrastructure/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis(config *config.RedisConfig) *redis.Client {
	host := config.Host
	username := config.Username
	password := config.Password
	port := config.Port
	db := config.Db

	Redis := redis.NewClient(&redis.Options{
		Username: username,
		Password: password,
		Addr:     host + ":" + port,
		DB:       db,
	})

	return Redis
}