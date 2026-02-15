package diff

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCache(addr string, ttl time.Duration) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Cache{client: rdb, ttl: ttl}
}

func (c *Cache) key(fromID, toID string) string {
	return fmt.Sprintf("diff:%s:%s", fromID, toID)
}

func (c *Cache) Get(ctx context.Context, fromID, toID string) (string, bool, error) {
	val, err := c.client.Get(ctx, c.key(fromID, toID)).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func (c *Cache) Set(ctx context.Context, fromID, toID, diff string) error {
	return c.client.Set(ctx, c.key(fromID, toID), diff, c.ttl).Err()
}
