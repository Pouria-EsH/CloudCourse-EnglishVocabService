package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	passw             string
	rdb               *redis.Client
	connectionTimeout time.Duration
	expTime           time.Duration
}

var ErrNotCached = errors.New("value not cached")

func NewRedisClient(address string, password string, cacheTimeout time.Duration, connectionTimeout time.Duration) *RedisCache {
	clnt := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	return &RedisCache{
		passw:             password,
		rdb:               clnt,
		connectionTimeout: connectionTimeout,
		expTime:           cacheTimeout,
	}
}

func (rc RedisCache) AddWord(word string, definition string) error {
	ctx, done := context.WithTimeout(context.Background(), rc.connectionTimeout)
	defer done()
	err := rc.rdb.Set(ctx, word, definition, rc.expTime).Err()
	return err
}

func (rc RedisCache) GetWord(word string) (string, error) {
	ctx, done := context.WithTimeout(context.Background(), rc.connectionTimeout)
	defer done()
	val, err := rc.rdb.Get(ctx, word).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrNotCached
		}
		return val, err
	}
	return val, nil
}
