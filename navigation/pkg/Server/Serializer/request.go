package Serializer

import (
	"github.com/labstack/echo/v4"
	"navigation/pkg/DTO"
)

func ConvertProtoToGetRouteRequest(ctx echo.Context) (DTO.NavigationRequest, error) {
	req := new(DTO.NavigationRequest)
	if err := ctx.Bind(req); err != nil {
		return DTO.NavigationRequest{}, err
	}
	return *req, nil
}
