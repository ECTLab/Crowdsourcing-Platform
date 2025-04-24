package redistool

import (
	"github.com/go-redis/redis/v8"
)

func (c *Client) Init(conf Config) {
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
		})
	} else {
		c.RedisClient = redis.NewClient(&redis.Options{
			Addr: addresses[0],
		})
	}
}
