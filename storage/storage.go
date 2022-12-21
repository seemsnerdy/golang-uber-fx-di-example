package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewRedisStorage)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
	ttl    time.Duration
}

func NewRedisStorage(
	client *redis.Client,
	ctx context.Context,
	ttl time.Duration,
) (*RedisStorage, error) {
	return &RedisStorage{
		client: client,
		ctx:    ctx,
		ttl:    ttl,
	}, nil
}

func (r *RedisStorage) GetCount(key string) (int64, error) {
	result, err := r.client.Get(r.ctx, key).Int64()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (r *RedisStorage) IncreaseCount(key string) error {
	if err := r.client.Incr(r.ctx, key).Err(); err != nil {
		return err
	}
	if err := r.client.Expire(r.ctx, key, r.ttl).Err(); err != nil {
		return err
	}
	return nil
}
