package Osrm

import (
	log "github.com/sirupsen/logrus"
	"navigation/config"
    "navigation/pkg/Clients/Http"
	"time"
)

func InitOsrmClients() {
	log.Info("trying to init osrm client")
	navigationOsrmNewClient = Client{
		Version:        "v1",
		Profile:        "driving",
		Timeout:        time.Duration(config.GetServiceConfig().Osrm.OsrmControl.TimeoutMS) * time.Millisecond,
		MaxRetries:     1,
		BackoffMax:     120 * time.Millisecond,
		BackoffFactor:  0.5,
		BackoffEnabled: false,
		Session:        Http.NewHttpClient(time.Duration(config.GetServiceConfig().Osrm.OsrmControl.TimeoutMS) * time.Millisecond),
		DecodeOutput:   false,
	}
	log.Info("osrm client initiated")
}



func InitNavigationOSRMService() {
	log.Info("trying to initiate osrm service")
	NavigationOsrmService = NavigationService{
		Name:                 "NavigationOsrmControl",
		Host:                 config.GetServiceConfig().Osrm.OsrmControl.Address,
		Client:               navigationOsrmNewClient,
		DefaultServiceClient: Http.NewHttpClient(time.Duration(config.GetServiceConfig().Osrm.OsrmControl.TimeoutMS) * time.Millisecond),
	}
	log.Info("osrm service initiated")
}
