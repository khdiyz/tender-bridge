package cache

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
		ctx:    context.Background(),
	}
}

func (c *RedisCache) Get(key string, dest any) error {
	data, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

func (c *RedisCache) Set(key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(c.ctx, key, data, ttl).Err()
}

func (c *RedisCache) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}
