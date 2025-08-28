package redisclient

import (
	"context"
	"fmt"
	"spitikos/api/internal/config"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
	cfg    *config.Config
}

func New(cfg *config.Config) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	return &Client{
		client: client,
		cfg:    cfg,
	}, nil
}

func (c *Client) Set(ctx context.Context, key string, value any) error {
	err := c.client.Set(ctx, key, value, 0).Err()
	return err
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	result, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s does not exist", key)
	}
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	var keys []string
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}
