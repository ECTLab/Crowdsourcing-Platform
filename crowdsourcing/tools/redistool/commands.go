package redistool

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func (c *Client) Get(key string) (string, error) {
	redisClient := c.RedisClient
	if redisClient != nil {
		ctx := context.Background()
		value, err := redisClient.Get(ctx, key).Result()
		if err == redis.Nil {
			return "", nil
		} else if err != nil {
			return "", err
		} else {
			return value, err
		}
	} else {
		return "", nil
	}
}

func (c *Client) Del(key string) error {
	redisClient := c.RedisClient
	if redisClient != nil {
		ctx := context.Background()
		err := redisClient.Del(ctx, key).Err()
		return err
	} else {
		return nil
	}
}

func (c *Client) Set(key string, value string) error {
	return c.SetWithTTL(key, value, 0)
}

func (c *Client) SetWithTTL(key string, value string, ttl time.Duration) error {
	redisClient := c.RedisClient
	if redisClient != nil {
		ctx := context.Background()
		err := redisClient.Set(ctx, key, value, ttl).Err()
		return err
	} else {
		return nil
	}
}

func (c *Client) HGet(hash string, key string) (string, error) {
	redisClient := c.RedisClient
	if redisClient != nil {
		ctx := context.Background()
		value, err := redisClient.HGet(ctx, hash, key).Result()
		if err == redis.Nil {
			return "", nil
		} else if err != nil {
			return "", err
		} else {
			return value, err
		}
	} else {
		return "", nil
	}
}

func (c *Client) HGetAll(hash string) (map[string]string, error) {
	redisClient := c.RedisClient
	var empty = map[string]string{}
	if redisClient != nil {
		ctx := context.Background()
		value, err := redisClient.HGetAll(ctx, hash).Result()
		if err == redis.Nil {
			return empty, nil
		} else if err != nil {
			return empty, err
		} else {
			return value, err
		}
	} else {
		return empty, nil
	}
}

func (c *Client) HSet(hash string, key string, value string) error {
	redisClient := c.RedisClient
	if redisClient != nil {
		ctx := context.Background()
		err := redisClient.HSet(ctx, hash, key, value).Err()
		return err
	} else {
		return nil
	}
}

func (c *Client) ListPush(key string, values ...string) error {
	redisClient := c.RedisClient
	if redisClient != nil {
		ctx := context.Background()
		err := redisClient.RPush(ctx, key, values).Err()
		return err
	} else {
		return nil
	}
}

func (c *Client) ListRangeAll(key string) ([]string, error) {
	redisClient := c.RedisClient
	if redisClient != nil {
		ctx := context.Background()
		result, err := redisClient.LRange(ctx, key, 0, -1).Result()
		if err == redis.Nil {
			return []string{}, nil
		}
		return result, err
	} else {
		return []string{}, nil
	}
}

func (c *Client) ListLPopN(key string, n int) error {
	redisClient := c.RedisClient
	if redisClient != nil {
		ctx := context.Background()
		err := redisClient.LPopCount(ctx, key, n).Err()
		if err == redis.Nil {
			return nil
		}
		return err
	} else {
		return nil
	}
}
