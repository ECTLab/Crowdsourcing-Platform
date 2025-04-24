package inridereports

import (
	"crowdsourcing/config"
	"crowdsourcing/pkg/clients/osrm"
	"crowdsourcing/tools/h3"
	"crowdsourcing/tools/osrmtool"
	"errors"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func createOnlineReport(report InRideReport) (string, error) {
	if !isOnlineReport(report.Type) {
		err := errors.New("bad online report type :" + string(report.Type))
		log.WithError(err).Errorf("create online report: bad online report type")
		return "", err
	}

	loc0 := report.EngagementLocationTime
	loc1 := report.SubmitLocationTime

	resp, err := osrm.Basic.Match([]osrmtool.Location{
		{loc0.Latitude, loc0.Longitude, 0, float64(loc0.Time), 0},
		{loc1.Latitude, loc1.Longitude, 0, float64(loc1.Time), 0},
	})
	if err != nil {
		log.WithError(err).WithField("report", report).Errorf("create online report: osrm match reponse with no tracepoints")
		return "", err
	}

	if len(resp.Tracepoints) == 0 {
		err := errors.New("OnlineReport osrm match with no tracepoint")
		log.WithError(err).WithFields(log.Fields{
			"osrm_response": resp,
			"report":        report,
		}).Errorf("create online report: osrm match reponse with no tracepoints")
		return "", err
	}

	location := resp.Tracepoints[0].Location
	seg := []int64{0, 0}
	if len(resp.Matchings) > 0 &&
		len(resp.Matchings[0].Legs) > 0 &&
		len(resp.Matchings[0].Legs[0].Annotation.Nodes) >= 2 {
		seg = resp.Matchings[0].Legs[0].Annotation.Nodes[0:2]
	}
	u := max(seg[0], seg[1])
	v := min(seg[0], seg[1])
	index := h3.GetH3Index(location[0], location[1])
	uid := uuid.NewString()
	timestamp := time.Now().UnixMilli()
	defaultTtlMins := config.GetServiceConfig().RedisCrowdsourcing.EventDefaultTtlMinutes

	event := Event{
		PositiveConfirmation: 1,
		NegativeConfirmation: 0,
		Type:                 OnlineReportType(report.Type),
		Location:             RedisLocation{location[1], location[0]}, // Osrm location is [longitude, latitude]
		U:                    u,
		V:                    v,
		H3Index:              index,
		Uid:                  uid,
		ReportTimestamp:      timestamp,
		UpdateTimestamp:      timestamp,
		TtlMinutes:           defaultTtlMins,
		IsAggregated:         false,
	}

	err = RedisAddNewOnlineReport(event)
	if err != nil {
		log.WithError(err).WithField("event", event).Errorf("create online report: could not add to redis")
		return "", err
	}

	reportId := redisEventKey{
		h3:        index,
		uid:       uid,
		eventType: OnlineReportType(report.Type),
	}.String()

	return reportId, nil
}

func isOnlineReport(reportType InRideReportType) bool {
	reportTypeIsOnline := false
	for _, t := range OnlineReportTypes {
		if OnlineReportType(reportType) == t {
			reportTypeIsOnline = true
		}
	}

	return reportTypeIsOnline
}
