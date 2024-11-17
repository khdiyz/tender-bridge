package cache

import (
	"encoding/json"
	"fmt"
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

func (c *RedisCache) DeletePattern(pattern string) error {
	iter := c.client.Scan(c.ctx, 0, pattern, 0).Iterator()
	for iter.Next(c.ctx) {
		if err := c.client.Del(c.ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("failed to delete key %s: %w", iter.Val(), err)
		}
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to iterate keys with pattern %s: %w", pattern, err)
	}

	return nil
}
