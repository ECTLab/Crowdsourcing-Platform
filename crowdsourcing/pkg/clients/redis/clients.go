package redis

import (
	"crowdsourcing/config"
	"crowdsourcing/tools/redistool"
)

var Crowdsourcing redistool.Client

func InitCrowdsourcing() {
	conf := config.GetServiceConfig().RedisCrowdsourcing
	rconf := redistool.Config{
		Enabled:      conf.Enabled,
		IsSentinel:   conf.IsSentinel,
		Addresses:    conf.Addresses,
		SentinelName: conf.SentinelName,
	}
	Crowdsourcing.Init(rconf)
}

var Annotation redistool.Client

func InitAnnotation() {
	conf := config.GetServiceConfig().RedisAnnotation
	rconf := redistool.Config{
		Enabled:      conf.Enabled,
		IsSentinel:   conf.IsSentinel,
		Addresses:    conf.Addresses,
		SentinelName: conf.SentinelName,
	}
	Annotation.Init(rconf)
}

func Init() {
	InitAnnotation()
	InitCrowdsourcing()
}
