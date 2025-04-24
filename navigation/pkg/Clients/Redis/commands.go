package Redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func (c *client) Get(key string) (string, error) {
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

func (c *client) Del(key string) error {
	redisClient := c.RedisClient
	if redisClient != nil {
		ctx := context.Background()
		err := redisClient.Del(ctx, key).Err()
		return err
	} else {
		return nil
	}
}

func (c *client) Set(key string, value string) error {
	return c.SetWithTTL(key, value, 0)
}

func (c *client) SetWithTTL(key string, value string, ttl time.Duration) error {
	redisClient := c.RedisClient
	if redisClient != nil {
		ctx := context.Background()
		err := redisClient.Set(ctx, key, value, ttl).Err()
		return err
	} else {
		return nil
	}
}

func (c *client) HGet(hash string, key string) (string, error) {
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

func (c *client) HGetAll(hash string) (map[string]string, error) {
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

func (c *client) HSet(hash string, key string, value string) error {
	redisClient := c.RedisClient
	if redisClient != nil {
		ctx := context.Background()
		err := redisClient.HSet(ctx, hash, key, value).Err()
		return err
	} else {
		return nil
	}
}

func (c *client) ListPush(key string, values ...string) error {
	redisClient := c.RedisClient
	if redisClient != nil {
		ctx := context.Background()
		err := redisClient.RPush(ctx, key, values).Err()
		return err
	} else {
		return nil
	}
}

func (c *client) ListRangeAll(key string) ([]string, error) {
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

func (c *client) ListLPopN(key string, n int) error {
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


func (c *client) SetPipeline(pairs map[string]string, ttl time.Duration) error {
	redisClient := c.RedisClient
	if redisClient != nil {
		ctx := context.Background()
		pipe := redisClient.Pipeline()
		for key, value := range pairs {
			pipe.Set(ctx, key, value, ttl)
		}

		_, err := pipe.Exec(ctx)


		return err
	} else {
		return nil
	}
}

func (c *client) GetPipeline (keys []string) (map[string]string, error) {
	redisClient := c.RedisClient
	if redisClient != nil {
		ctx := context.Background()
		pipe := redisClient.Pipeline()
		cmds := make([]*redis.StringCmd, len(keys))
		for i, key := range keys {
			cmds[i] = pipe.Get(ctx, key)
		}

		_, err := pipe.Exec(ctx)
		if err != nil && err != redis.Nil {
			return nil, err
		}

		results := make(map[string]string)
		for i, cmd := range cmds {
			value, err := cmd.Result()
			if err == redis.Nil {
				results[keys[i]] = ""
			} else if err != nil {
				return nil, err
			} else {
				results[keys[i]] = value
			}
		}

		return results, nil
	} else {
		return nil, nil
	}
}
