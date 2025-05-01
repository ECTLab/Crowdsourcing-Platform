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
	}
	response := inridereports.SerializeConfirmResponse(nil)

	return ctx.JSON(http.StatusAccepted, response)
}

func GetVisibilityReport(ctx echo.Context) error {

	type visibilityItem [4]any
	type visibilityReport []visibilityItem

	reportType := inridereports.POLICE_ONLINE

	allOnlineEvents, err := inridereports.RedisReadAllOnlineReports()
	if err != nil {
		log.WithError(err).Errorf("error reading all crowdsourcing reports from redis hash")
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	var reports = visibilityReport{}
	for _, events := range allOnlineEvents[reportType] {
		for _, event := range events {
			yes := float64(event.PositiveConfirmation)
			no := float64(event.NegativeConfirmation)
			confidence := float64((yes * 3) / ((yes * 3) + no))

			report := visibilityItem{
				event.Location[0],
				event.Location[1],
				event.IsAggregated,
				confidence,
			}

			reports = append(reports, report)
		}
	}
	return ctx.JSON(http.StatusOK, reports)
}
