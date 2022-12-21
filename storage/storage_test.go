package storage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
)

func TestRedisStorage_GetCount(t *testing.T) {
	redisClient, redisClientMock := redismock.NewClientMock()
	type fields struct {
		client *redis.Client
		ctx    context.Context
		ttl    time.Duration
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
		mock    func()
	}{
		{
			name: "simple",
			fields: fields{
				client: redisClient,
				ctx:    context.TODO(),
				ttl:    time.Minute,
			},
			args:    args{key: "random key"},
			want:    10,
			wantErr: false,
			mock: func() {
				redisClientMock.ExpectGet("random key").SetVal("10")
			},
		},
		{
			name: "error",
			fields: fields{
				client: redisClient,
				ctx:    context.TODO(),
				ttl:    time.Minute,
			},
			args:    args{key: "random key"},
			want:    0,
			wantErr: true,
			mock: func() {
				redisClientMock.ExpectGet("random key").SetErr(fmt.Errorf("some error"))
			},
		},
		{
			name: "nil",
			fields: fields{
				client: redisClient,
				ctx:    context.TODO(),
				ttl:    time.Minute,
			},
			args:    args{key: "random key"},
			want:    0,
			wantErr: true,
			mock: func() {
				redisClientMock.ExpectGet("random key").RedisNil()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			r := &RedisStorage{
				client: tt.fields.client,
				ctx:    tt.fields.ctx,
				ttl:    tt.fields.ttl,
			}
			got, err := r.GetCount(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCount() got = %v, want %v", got, tt.want)
			}
			if err := redisClientMock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestRedisStorage_IncreaseCount(t *testing.T) {
	redisClient, redisClientMock := redismock.NewClientMock()
	type fields struct {
		client *redis.Client
		ctx    context.Context
		ttl    time.Duration
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name: "simple",
			fields: fields{
				client: redisClient,
				ctx:    context.TODO(),
				ttl:    time.Minute,
			},
			args:    args{key: "test key"},
			wantErr: false,
			mock: func() {
				redisClientMock.ExpectIncr("test key").SetVal(1)
				redisClientMock.ExpectExpire("test key", time.Minute).SetVal(true)
			},
		},
		{
			name: "error on expire",
			fields: fields{
				client: redisClient,
				ctx:    context.TODO(),
				ttl:    time.Minute,
			},
			args:    args{key: "test key"},
			wantErr: true,
			mock: func() {
				redisClientMock.ExpectIncr("test key").SetVal(1)
				redisClientMock.ExpectExpire("test key", time.Minute).SetErr(fmt.Errorf("some error"))
			},
		},
		{
			name: "error on incr",
			fields: fields{
				client: redisClient,
				ctx:    context.TODO(),
				ttl:    time.Minute,
			},
			args:    args{key: "test key"},
			wantErr: true,
			mock: func() {
				redisClientMock.ExpectIncr("test key").SetErr(fmt.Errorf("some error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			r := &RedisStorage{
				client: tt.fields.client,
				ctx:    tt.fields.ctx,
				ttl:    tt.fields.ttl,
			}
			if err := r.IncreaseCount(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("IncreaseCount() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := redisClientMock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
