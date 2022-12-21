package configuration

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewRedisConfigurationFromEnv)

type RedisConfiguration struct {
	Address string `env:"REDIS_ADDRESS,required"`
}

func NewRedisConfigurationFromEnv() (*RedisConfiguration, error) {
	var result RedisConfiguration
	if err := env.Parse(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *RedisConfiguration) Options() *redis.Options {
	return &redis.Options{Addr: r.Address}
}
