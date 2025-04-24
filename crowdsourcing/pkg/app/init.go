package app

import (
	"crowdsourcing/config"
	"crowdsourcing/pkg/clients/osrm"
	"crowdsourcing/pkg/clients/redis"
	stdlog "log"
	"os"
)

func Init() {
	config.InitViper()
	stdlog.SetOutput(os.Stdout)
	stdlog.Printf("Starting %s service", config.GetServiceConfig().ServiceConfigs.ServiceName)




	redis.Init()
	osrm.Init()

}
