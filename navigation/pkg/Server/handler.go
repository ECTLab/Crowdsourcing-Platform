package Server

import (
    "fmt"
    "github.com/labstack/echo/v4"
    "navigation/pkg/Route"
    "navigation/pkg/Server/Serializer"
    "net/http"

    log "github.com/sirupsen/logrus"
)

func GetRoute(ctx echo.Context) error {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("recovered panic: %v", r)
		}
	}()

	log.Info("Request accepted")

	structRequest, err := Serializer.ConvertProtoToGetRouteRequest(ctx)
	if err != nil {
		log.Errorf("error in getting request: %s", err.Error())
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	routeResponse, err := Route.GetRoute(&structRequest)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to provide voice instructions:%s", err.Error()))
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	log.Info("Request Handled successfully")
	return ctx.JSON(http.StatusOK, routeResponse)

}
