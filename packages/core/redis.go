package core

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	RedisKeyTtlTime = 5 * time.Minute
)

func NewRedisClient(redisURI string) (*redis.Client, error) {
	// Parse Redis URI and create options
	opts, err := redis.ParseURL(redisURI)
	if err != nil {
		return nil, err
	}

	// Create Redis client
	client := redis.NewClient(opts)

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

type CacheRedis struct {
	RedisClient *redis.Client
}

type ICache interface {
	CloseConnection()
	Exists(key string) (bool, error)
	Get(key string) (string, error)
	Set(key string, value string) (bool, error)
}

func NewCacheRedis(redisURI string) *CacheRedis {
	client, err := NewRedisClient(redisURI)
	if err != nil {
		return nil
	}

	return &CacheRedis{RedisClient: client}
}

func (cache *CacheRedis) CloseConnection() {
	_ = cache.RedisClient.Close()
}

func (cache *CacheRedis) Exists(key string) (bool, error) {
	exist, err := cache.RedisClient.GetEx(context.Background(), key, RedisKeyTtlTime).Bool()
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return exist, nil
}

func (cache *CacheRedis) Set(key string, value string) (bool, error) {
	return cache.SetEx(key, value, RedisKeyTtlTime)
}

func (cache *CacheRedis) SetEx(key string, value string, ttl time.Duration) (bool, error) {
	result, err := cache.RedisClient.SetEx(context.Background(), key, value, ttl).Result()
	if err != nil {
		return false, err
	}
	return result == "OK", nil
}

func (cache *CacheRedis) Get(key string) (string, error) {
	strResult, err := cache.RedisClient.GetEx(context.Background(), key, RedisKeyTtlTime).Result()
	if err != nil {
		return "", err
	}
	return strResult, nil
}
