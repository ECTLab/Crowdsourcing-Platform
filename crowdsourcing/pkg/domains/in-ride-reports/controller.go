package inridereports


func CreateReport(report InRideReport) (string, error) {
	if report.Online {
		return createOnlineReport(report)
	}

	return "", nil
}
