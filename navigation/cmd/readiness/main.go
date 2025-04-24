package main

import (
	"fmt"
	"navigation/config"
	health "navigation/tools/liveness"
)

func main() {
	config.InitViper()
	servicePort := config.GetServiceConfig().ServiceConfigs.ServiceProtoPort
	address := fmt.Sprintf("localhost:%d", servicePort)
	health.PerformHealthCheckRequest(address)
}
