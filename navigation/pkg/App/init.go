package App

import (
	"context"
	stdlog "log"
	"navigation/config"
	"navigation/pkg/Annotation"
	osrm "navigation/pkg/Clients/Osrm"
	"navigation/pkg/Clients/Redis"
	"navigation/tools/logger"
	"os"

)

func Init() {
	config.InitViper()
	stdlog.SetOutput(os.Stdout)
	stdlog.Printf("Starting %s service", config.GetServiceConfig().ServiceConfigs.Name)

	logger.Init(
		logger.Config{
			ServiceName: config.GetServiceConfig().ServiceConfigs.Name,
			LogLevel:    config.GetServiceConfig().LogLevel,
		},
	)

	Redis.Init()
	osrm.InitOsrmClients()
	osrm.InitNavigationOSRMService()
	if config.GetServiceConfig().RedisAnnotation.Enabled {
		if config.GetServiceConfig().Annotation.OnlineReportsEnabled {
			Annotation.UpdateOnlineReports(context.Background())
		}
	}
}
