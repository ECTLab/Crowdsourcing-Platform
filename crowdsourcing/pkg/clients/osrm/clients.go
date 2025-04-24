package osrm

import (
	"crowdsourcing/config"
	"crowdsourcing/tools/osrmtool"
)

var Basic osrmtool.Client

func InitBasic() {
	conf := config.GetServiceConfig().OsrmTable
	osrm_conf := osrmtool.Config{
		Address:   conf.Address,
		TimeoutMS: conf.TimeoutMS,
	}
	Basic.Init(osrm_conf)
}

var Matching osrmtool.Client

func InitMatching() {
	conf := config.GetServiceConfig().OsrmMatching
	osrm_conf := osrmtool.Config{
		Address:   conf.Address,
		TimeoutMS: conf.TimeoutMS,
	}
	Matching.Init(osrm_conf)
}

func Init() {
	InitBasic()
	InitMatching()
}
