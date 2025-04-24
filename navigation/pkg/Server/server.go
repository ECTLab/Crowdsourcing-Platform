package Server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"navigation/config"
	"navigation/pkg/Middleware"
)

func ServeHttp() {
		port := config.GetServiceConfig().ServiceConfigs.ServiceProtoPort
	e := echo.New()
	e.Use(Middleware.SessionCheckMiddleware())
	e.POST("/navigation/get-route", GetRoute)
	log.Info(fmt.Sprintf("serving on port %d", port))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
