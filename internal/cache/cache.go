package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct{ rdb *redis.Client }

func New(addr string) *Cache {
	return &Cache{rdb: redis.NewClient(&redis.Options{Addr: addr, DB: 0})}
}

func (c *Cache) Close() error { return c.rdb.Close() }

func (c *Cache) Set(key, value string, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return c.rdb.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return c.rdb.Get(ctx, key).Result()
}

func (c *Cache) TTL(key string) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return c.rdb.TTL(ctx, key).Result()
}
