package redis_client

import (
	"di-go-example/configuration"

	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewClient)

func NewClient(configuration *configuration.RedisConfiguration) *redis.Client {
	return redis.NewClient(configuration.Options())
}
