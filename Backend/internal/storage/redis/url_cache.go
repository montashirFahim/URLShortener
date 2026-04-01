package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type urlCache struct {
	client *redis.Client
}

func NewURLCache(client *redis.Client) *urlCache {
	return &urlCache{client: client}
}

func (c *urlCache) GetLongURL(ctx context.Context, code string) (string, error) {
	val, err := c.client.Get(ctx, "url:"+code).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

func (c *urlCache) SetLongURL(ctx context.Context, code string, longURL string, expiresAt *string) error {
	expiration := 24 * time.Hour
	if expiresAt != nil {
		t, err := time.Parse(time.RFC3339, *expiresAt)
		if err == nil {
			ttl := time.Until(t)
			if ttl < expiration {
				expiration = ttl
			}
		}
	}
	if expiration <= 0 {
		return nil
	}
	return c.client.Set(ctx, "url:"+code, longURL, expiration).Err()
}

func (c *urlCache) DeleteURL(ctx context.Context, code string) error {
	return c.client.Del(ctx, "url:"+code).Err()
}
