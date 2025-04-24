package inridereports

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"time"
)

func ConfirmPoliceReport(confirmation ReportConfirmation) error {
	reportId := confirmation.ReportId
	yes := confirmation.Confirmed

	var event Event
	event, err := RedisGetOnlineReport(reportId)
	if err != nil {
		return err
	}

	if !event.IsAggregated {
		err := errors.New("request for not aggregated report")
		log.WithError(err).WithField("event", event).Errorf("confirm online report: request for not aggregated report")
		return nil
	}

	if yes {
		event.PositiveConfirmation += 1
	} else {
		event.NegativeConfirmation += 1
	}

	timestamp := time.Now().UnixMilli()
	event.UpdateTimestamp = timestamp

	err = RedisSetOnlineReport(reportId, event)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"event":        event,
			"reportId":     reportId,
			"confirmation": confirmation,
		}).Errorf("confirm online report: could not update redis")
		return err
	}

	_ = confirmation.Reporter

	return nil
}

func ConfirmFeasibilityReport(confirmation ReportConfirmation) error {
	report, err := RedisGetFeasibilityConfirmationCase(confirmation.ReportId)
	if err != nil {
		log.WithFields(log.Fields{
			"confirmation":report,
			"request_details": confirmation,
		}).Errorf("failed to fetch annotation redis in feasibility confirmation: %s", err.Error())
		return err
	}
	report.Contributions += 1
	if err != nil {
		log.WithFields(log.Fields{
			"confirmation":report,
			"request_details": confirmation,
		}).Errorf("failed to stream feasibility confirmation: %s", err.Error())
	}
	err = RedisSetFeasibilityConfirmationCase(confirmation.ReportId, report)
	if err != nil {
		log.WithFields(log.Fields{
			"confirmation":report,
			"request_details": confirmation,
		}).Errorf("failed to update annotation redis in feasibility confirmation: %s", err.Error())
	}
	return err
}
