package repositories

import (
	"context"

	"github.com/go-redis/redis/v8"
)

const readArticleIDListKey = "items.read_dev_article_ids"

type DevArticleRepository interface {
	ListAllReadArticleID(ctx context.Context) ([]string, error)
	AddReadArticleID(ctx context.Context, id string) error
}

type RedisDevArticleRepository struct {
	BaseRedisRepository
}

func NewRedisDevArticleRepository(rdb *redis.Client) (DevArticleRepository, error) {
	repo := &RedisDevArticleRepository{}
	repo.rdb = rdb
	return repo, nil
}

func (repo *RedisDevArticleRepository) ListAllReadArticleID(ctx context.Context) ([]string, error) {
	return repo.rdb.LRange(ctx, readArticleIDListKey, 0, -1).Result()
}

func (repo *RedisDevArticleRepository) AddReadArticleID(ctx context.Context, id string) error {
	err := repo.rdb.LPush(ctx, readArticleIDListKey, id).Err()
	if err != nil {
		return err
	}

	return repo.rdb.LTrim(ctx, readArticleIDListKey, 0, 13).Err()
}
