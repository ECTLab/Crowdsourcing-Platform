package Redis

import (
	"github.com/go-redis/redis/v8"
	"time"
)


type client struct {
	RedisClient redis.UniversalClient
}

type config struct {
	Enabled      bool
	IsSentinel   bool
	Addresses    []string
	SentinelName string
}

func (c *client) Init(conf config) {
	enabled := conf.Enabled
	if !enabled {
		return
	}

	addresses := conf.Addresses
	isSentinel := conf.IsSentinel
	sentinelName := conf.SentinelName
	if isSentinel {
		c.RedisClient = redis.NewUniversalClient(&redis.UniversalOptions{
			MasterName: sentinelName,
			Addrs:      addresses,
			DialTimeout: 2 * time.Second,
		})
	} else {
		c.RedisClient = redis.NewClient(&redis.Options{
			Addr: addresses[0],
		})
	}
}