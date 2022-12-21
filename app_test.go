package main

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"go.uber.org/fx"
)

func TestNewApp(t *testing.T) {
	redisClient, redisClientMock := redismock.NewClientMock()
	tests := []struct {
		name string
		mock func()
	}{
		{
			name: "simple",
			mock: func() {
				redisClientMock.ExpectIncr("some key").SetVal(10)
				redisClientMock.ExpectExpire("some key", 2*time.Minute).SetVal(true)
				redisClientMock.ExpectGet("some key").SetVal("10")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			app := fx.New(
				NewApp(),
				fx.Replace(
					redisClient,
				),
				fx.Replace(
					time.Minute*2,
				),
			)
			if err := app.Start(context.TODO()); err != nil {
				t.Error(err)
			}
			if err := redisClientMock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
