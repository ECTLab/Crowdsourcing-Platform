package server

import (
	"crowdsourcing/config"
	"crowdsourcing/pkg/middleware"
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	_ "google.golang.org/grpc/encoding/gzip"
)

func ServeHttp() {
	port := config.GetServiceConfig().ServiceConfigs.ServiceProtoPort
	e := echo.New()
	e.Use(middleware.SessionCheckMiddleware())
	e.POST("/crowdsourcing/in-ride-report", CreateInRideReport)
	e.POST("/crowdsourcing/in-ride-report/:report_id/confirm", ConfirmInRideReport)
	e.GET("/crowdsourcing/reports/visibility", GetVisibilityReport)
	log.Info(fmt.Sprintf("serving on port %d", port))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
