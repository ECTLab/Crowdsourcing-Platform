package inridereports

import (
	"github.com/labstack/echo/v4"
)

func DeserializeRequest(ctx echo.Context) (InRideReport, error) {

	req := new(InRideReport)
	if err := ctx.Bind(req); err != nil {
		return InRideReport{}, err
	}
	if isReportOnline(req.Type) {
		req.Online = true
	}
	return *req, nil
}

func SerializeResponse(reportId string) InRideReportResponse {
	return InRideReportResponse{
		ReportId: reportId,
	}
}



func DeserializeConfirmRequest(ctx echo.Context) (ReportConfirmation, error) {
	req := new(ReportConfirmation)
	if err := ctx.Bind(req); err != nil {
		return ReportConfirmation{}, err
	}
	req.ReportId = ctx.Param("report_id")
	return *req, nil
}

func SerializeConfirmResponse(response any) ReportConfirmationResponse {
	return ReportConfirmationResponse{}
}

func isReportOnline(reportType InRideReportType) bool {
	if reportType == InRideReportType(POLICE_ONLINE) || reportType == InRideReportType(ACCIDENT_ONLINE) {
		return true
	}
	return false
}