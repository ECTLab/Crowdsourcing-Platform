package Redis

import (
	log "github.com/sirupsen/logrus"
	serviceConfigs "navigation/config"
)


var Crowdsourcing client

func InitCrowdsourcing() {
	conf := serviceConfigs.GetServiceConfig().RedisAnnotation
	rconf := config{
		Enabled:      conf.Enabled,
		IsSentinel:   conf.IsSentinel,
		Addresses:    conf.Addresses,
		SentinelName: conf.SentinelName,
	}
	Crowdsourcing.Init(rconf)
	log.Info("crowdsource redis initiated")
}

var Traffic client

func InitTraffic() {
	conf := serviceConfigs.GetServiceConfig().RedisAnnotation
	rconf := config{
		Enabled:      conf.Enabled,
		IsSentinel:   conf.IsSentinel,
		Addresses:    conf.Addresses,
		SentinelName: conf.SentinelName,
	}
	Traffic.Init(rconf)
	log.Info("traffic redis initiated")
}

var Annotation client

func InitAnnotation() {
	conf := serviceConfigs.GetServiceConfig().RedisAnnotation
	rconf := config{
		Enabled:      conf.Enabled,
		IsSentinel:   conf.IsSentinel,
		Addresses:    conf.Addresses,
		SentinelName: conf.SentinelName,
	}
	Annotation.Init(rconf)
	log.Info("annotation redis initiated")
}

var TTS client

func InitTTS() {
	conf := serviceConfigs.GetServiceConfig().RedisTTS
	rconf := config{
		Enabled:      conf.Enabled,
		IsSentinel:   conf.IsSentinel,
		Addresses:    conf.Addresses,
		SentinelName: conf.SentinelName,
	}
	TTS.Init(rconf)
	log.Info("tts redis initiated")
}

var Feasibility client

func InitFeasibility() {
	conf := serviceConfigs.GetServiceConfig().RedisAnnotation
	rconf := config {
		Enabled: conf.Enabled,
		IsSentinel: conf.IsSentinel,
		Addresses: conf.Addresses,
		SentinelName: conf.SentinelName,
	}
	Feasibility.Init(rconf)
	log.Info("feasibility redis initiated")
}

func Init() {
	InitAnnotation()
	InitCrowdsourcing()
	InitTraffic()
	InitTTS()
	InitFeasibility()
}
