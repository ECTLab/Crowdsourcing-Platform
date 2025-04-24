package redistool

import (
	"github.com/go-redis/redis/v8"
)

type Client struct {
	RedisClient redis.UniversalClient
}

type Config struct {
	Enabled      bool
	IsSentinel   bool
	Addresses    []string
	SentinelName string
}
