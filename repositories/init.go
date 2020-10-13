package repositories

import (
	"github.com/go-redis/redis/v8"
)

type BaseRedisRepository struct {
	rdb *redis.Client
}
