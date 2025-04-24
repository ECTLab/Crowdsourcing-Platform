package server

import (
	inridereports "crowdsourcing/pkg/domains/in-ride-reports"
	"github.com/labstack/echo/v4"
	"net/http"


	log "github.com/sirupsen/logrus"
)


func CreateInRideReport(ctx echo.Context) error {
	request, err := inridereports.DeserializeRequest(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	uid, err := inridereports.CreateReport(request)
	if err != nil {
		log.WithError(err).Errorf("Error creating Report : %v", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	response := inridereports.SerializeResponse(uid)

	return ctx.JSON(http.StatusCreated, response)
}

func ConfirmInRideReport(ctx echo.Context) error {
	request, err := inridereports.DeserializeConfirmRequest(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	if request.Type == inridereports.POLICE {
		err := inridereports.ConfirmPoliceReport(request)
		if err != nil {
			log.WithError(err).Errorf("Error confirming Police Report : %v", err)
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
	} else if request.Type == inridereports.FEASIBILITY {
		err := inridereports.ConfirmFeasibilityReport(request)
		if err != nil {
			log.WithError(err).Errorf("Error confirming Feasibility Report : %v", err)
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
	}

	response := inridereports.SerializeConfirmResponse(nil)

	return ctx.JSON(http.StatusAccepted, response)
}
