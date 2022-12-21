package main

import (
	"context"
	"log"
	"time"

	"di-go-example/configuration"
	"di-go-example/redis_client"
	"di-go-example/storage"

	"go.uber.org/fx"
)

func NewApp() fx.Option {
	return fx.Options(
		configuration.Module,
		storage.Module,
		redis_client.Module,
		fx.Provide(context.Background),
		fx.Supply(time.Hour*2),
		fx.Invoke(Runnable),
	)
}

func Runnable(s *storage.RedisStorage) {
	key := "some key"
	if err := s.IncreaseCount(key); err != nil {
		panic(err)
	}
	result, err := s.GetCount(key)
	if err != nil {
		panic(err)
	}
	log.Printf("result is: %v\n", result)
}
