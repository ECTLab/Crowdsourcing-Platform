package inridereports

import (
	"crowdsourcing/config"
	"crowdsourcing/pkg/clients/osrm"
	"crowdsourcing/tools/clustering"
	"crowdsourcing/tools/h3"
	"crowdsourcing/tools/osm"
	"crowdsourcing/tools/osrmtool"
	"math"
	"math/rand"
	"time"

	"github.com/qedus/osmpbf"
	log "github.com/sirupsen/logrus"
)

func AggregateAllOnlineReports() error {
	startTime := time.Now().UnixMilli()

	ClusterDurationThreshold := config.GetServiceConfig().AggregationConfigs.ClusterDurationThreshold
	ClusterLouvainResolution := config.GetServiceConfig().AggregationConfigs.ClusterLouvainResolution
	MaxOsrmTableLocations := config.GetServiceConfig().AggregationConfigs.MaxOsrmTableLocations

	log.Infof(
		"starting aggregation with config: {ClusterDurationThreshold:%.4f ClusterLouvainResolution:%.4f MaxOsrmTableLocations: %d}",
		ClusterDurationThreshold,
		ClusterLouvainResolution,
		MaxOsrmTableLocations,
	)

	downloadOsmPbf()

	var allReports map[OnlineReportType]map[string][]Event
	var newReports []redisEventKey

	log.Infof("reading all reports from redis")

	allReports, newReports, err := RedisReadAllNewOnlineReports()
	if err != nil {
		return err
	}


	log.Infof("aggregating events")

	var allAggregatedEvents map[string]Event // uid -> event
	var newReportUids map[string]struct{}    // set(uid)
	var doneReports map[string]struct{}      // set(uid)

	allAggregatedEvents, newReportUids, doneReports, err = aggregateEvents(
		allReports,
		newReports,
		startTime,
		ClusterDurationThreshold,
		ClusterLouvainResolution,
		MaxOsrmTableLocations,
	)
	if err != nil {
		return err
	}

	log.Infof("updating redis with aggregated events")

	err = updateAggregatedEvents(allAggregatedEvents, newReportUids, doneReports)
	if err != nil {
		return err
	}

	return nil
}

func aggregateEvents(
	allReports map[OnlineReportType]map[string][]Event,
	newReports []redisEventKey,
	startTime int64,
	ClusterDurationThreshold float64,
	ClusterLouvainResolution float64,
	MaxOsrmTableLocations int,
) (
	map[string]Event,
	map[string]struct{},
	map[string]struct{},
	error,
) {
	var allAggregatedEvents = map[string]Event{} // uid -> event
	var newReportUids = map[string]struct{}{}    // set(uid)
	var doneReports = map[string]struct{}{}      // set(uid)

	for _, newReport := range newReports {
		newReportUids[newReport.uid] = struct{}{}
		_, ok := doneReports[newReport.uid]
		if ok {
			continue
		}

		var eventsToAgg []Event
		var eventLocations []osrmtool.Location

		eventsToAgg, eventLocations, err := extractUsingH3WithNeighbours(newReport, allReports, doneReports, startTime)
		if err != nil {
			return nil, nil, nil, err
		}

		removeExcessEvents(&eventLocations, &eventsToAgg, MaxOsrmTableLocations)

		clusters, centers, err := clusterEvents(eventLocations, ClusterDurationThreshold, ClusterLouvainResolution)
		if err != nil {
			return nil, nil, nil, err
		}

		extractEventsFromClusters(allAggregatedEvents, eventsToAgg, clusters, centers, startTime)
	}

	return allAggregatedEvents, newReportUids, doneReports, nil
}

func extractUsingH3WithNeighbours(
	report redisEventKey,
	allReports map[OnlineReportType]map[string][]Event,
	doneReports map[string]struct{},
	startTime int64,
) (
	[]Event,
	[]osrmtool.Location,
	error,
) {
	var eventsToAgg = []Event{}
	var eventLocations = []osrmtool.Location{}

	reportCell, err := h3.GetH3CellFromStringIndex(report.h3)
	if err != nil {
		log.WithError(err).Errorf("aggregation: error: decoding h3")
		return nil, nil, err
	}
	h3WithNeighbours := h3.AsIndex(h3.WithNeighbours(reportCell))

	for _, h3 := range h3WithNeighbours {
		events := allReports[report.eventType][h3]
		for _, event := range events {
			doneReports[event.Uid] = struct{}{}
			if isReportExpired(event, startTime) {
				continue
			}
			eventsToAgg = append(eventsToAgg, event)
			osrmLoc := osrmtool.Location{event.Location[0], event.Location[1]}
			eventLocations = append(eventLocations, osrmLoc)
		}
	}

	return eventsToAgg, eventLocations, nil
}

func removeExcessEvents(locationsPtr *[]osrmtool.Location, eventsPtr *[]Event, maxAllowed int) {
	var locations = *locationsPtr
	var events = *eventsPtr

	// randomly remove locations to fit in MaxOsrmTableLocations
	nLocs := len(locations)
	if nLocs > maxAllowed {
		rem := make([]bool, nLocs)
		newLocsSlice := []osrmtool.Location{}
		newEventsSlice := []Event{}

		// generate nLocs random numbers as indexes of locations we need to remove
		nToRem := nLocs - maxAllowed
		for i := 0; i < nToRem; i++ {
			idx := rand.Intn(nLocs)
			for rem[idx] {
				idx = rand.Intn(nLocs)
			}
			rem[idx] = true
		}

		for i := range locations {
			if !rem[i] {
				newLocsSlice = append(newLocsSlice, locations[i])
				newEventsSlice = append(newEventsSlice, events[i])
			}
		}
		*locationsPtr = newLocsSlice
		*eventsPtr = newEventsSlice
	}
}

func clusterEvents(locations []osrmtool.Location, durationThreshold, louvainResolution float64) ([][]int64, []int64, error) {
	var durations [][]float64 // from OSRM table
	var clusters [][]int64
	var centers []int64

	if len(locations) == 1 {
		clusters = [][]int64{{0}}
		centers = []int64{0}
	} else if len(locations) > 0 {
		res, err := osrm.Basic.Table(locations, nil, nil)
		if err != nil {
			log.WithError(err).Errorf("aggregation: error: osrm table")
			return nil, nil, err
		}

		durations = res.Durations

		var rows = []int64{}
		var cols = []int64{}
		for i := 0; i < len(durations); i++ {
			rows = append(rows, int64(i))
			cols = append(cols, int64(i))
		}

		clusters, centers = clustering.FindClusters(rows, cols, durations, durationThreshold, louvainResolution)
	}

	return clusters, centers, nil
}
func extractEventsFromClusters(clusteredEvents map[string]Event, events []Event, clusters [][]int64, centers []int64, startTime int64) {
	// extract aggregated events from clusters
	for i, cluster := range clusters {
		clusterCenterEvent := events[centers[i]]
		yes := 0
		no := 0
		uid := clusterCenterEvent.Uid
		for _, eventIndex := range cluster {
			event := events[eventIndex]
			if event.IsAggregated {
				// enable to not update the previous center report
				// clusterCenterEvent = event
				clusterCenterEvent.Uid = uid
			}
			yes += int(event.PositiveConfirmation)
			no += int(event.NegativeConfirmation)
		}

		clusterCenterEvent.PositiveConfirmation = int32(yes)
		clusterCenterEvent.NegativeConfirmation = int32(no)
		clusterCenterEvent.UpdateTimestamp = startTime
		clusterCenterEvent.IsAggregated = true
		clusteredEvents[clusterCenterEvent.Uid] = clusterCenterEvent
	}
}

func updateAggregatedEvents(aggregated map[string]Event, newReportUids map[string]struct{}, seenReports map[string]struct{}) error {
	var newReportsSlice = map[OnlineReportType]map[string][]Event{
		POLICE_ONLINE:   {},
		ACCIDENT_ONLINE: {},
	}
	allReports, err := RedisReadAllOnlineReports()
	if err != nil {
		log.WithError(err).Errorf("aggregation: error: reading redis at update phase")
		return err
	}

	now := time.Now().UnixMilli()

	// reportsList := []Event{}

	for reportType, reportsByH3 := range allReports {
		for h3index, reports := range reportsByH3 {
			for _, report := range reports {

				_, checkedDuringThisAggregation := newReportUids[report.Uid]
				aggregatedBefore := report.IsAggregated
				aggregatedReport, isResultOfThisAggregation := aggregated[report.Uid]
				_, participatedInThisAggregation := seenReports[report.Uid]

				submittedDuringThisAggregation := !checkedDuringThisAggregation && !aggregatedBefore && !participatedInThisAggregation

				Expired := isReportExpired(report, now)

				mustRemainInState := (submittedDuringThisAggregation ||
					isResultOfThisAggregation ||
					(aggregatedBefore && !participatedInThisAggregation)) && !Expired

				if mustRemainInState {
					reportToSave := report
					if isResultOfThisAggregation {
						reportToSave = aggregatedReport
					}
					newReportsSlice[reportType][h3index] = append(newReportsSlice[reportType][h3index], reportToSave)
				}

			}
		}
	}

	err = RedisReplaceAllOnlineReports(newReportsSlice)
	if err != nil {
		log.WithError(err).Errorf("aggregation: error: writing to redis at update phase")
		return err
	}

	var annotationReports = []*AnnotationOnlineReport{}
	var annotationReportsByU = map[int64][]*AnnotationOnlineReport{}
	var annotationReportsByV = map[int64][]*AnnotationOnlineReport{}
	for _, reportsByH3 := range newReportsSlice {
		for _, reports := range reportsByH3 {
			for _, report := range reports {
				if report.IsAggregated {
					yes := report.PositiveConfirmation
					no := report.NegativeConfirmation
					annotReport := AnnotationOnlineReport{
						Confidence:   float64(yes*3) / float64(no+(yes*3)),
						Confirmation: true,
						Type:         report.Type,
						U:            max(report.U, report.V),
						V:            min(report.U, report.V),
						Offset:       0.0,
						OffsetRest:   0.0,
						Uid: redisEventKey{
							h3:        report.H3Index,
							eventType: report.Type,
							uid:       report.Uid,
						}.String(),
						location: report.Location,
					}
					annotationReports = append(annotationReports, &annotReport)

					if _, ok := annotationReportsByU[annotReport.U]; !ok {
						annotationReportsByU[annotReport.U] = []*AnnotationOnlineReport{}
					}
					annotationReportsByU[annotReport.U] = append(annotationReportsByU[annotReport.U], &annotReport)

					if _, ok := annotationReportsByV[annotReport.V]; !ok {
						annotationReportsByV[annotReport.V] = []*AnnotationOnlineReport{}
					}
					annotationReportsByV[annotReport.V] = append(annotationReportsByV[annotReport.V], &annotReport)
				}
			}
		}
	}

	err = updateOffsets(annotationReportsByU, annotationReportsByV)
	if err != nil {
		log.WithError(err).Errorf("aggregation: error: calculating offsets")
		return err
	}

	err = RedisAnnotationUpdateReports(annotationReports)
	if err != nil {
		log.WithError(err).Errorf("aggregation: error: updating annotation redis")
		return err
	}

	return nil
}

func isReportExpired(report Event, atTimestamp int64) bool {
	return report.UpdateTimestamp+int64(report.TtlMinutes*60*1000) < atTimestamp
}



func updateOffsets(
	reportsByU map[int64][]*AnnotationOnlineReport,
	reportsByV map[int64][]*AnnotationOnlineReport,
) error {
	QuickDistanceCalc := func(lat1, long1, lat2, long2 float64) float64 {
		x := lat2 - lat1
		y := (long2 - long1) * math.Cos((lat2+lat1)*0.00872664626)
		distanceMeters := (111.319 * math.Sqrt(x*x+y*y)) * 1000
		return distanceMeters
	}
	OsmNodeIter := func(node any) error {
		switch node := node.(type) {
		case *osmpbf.Node:
			lat1 := node.Lat
			long1 := node.Lon

			reports, ok := reportsByU[node.ID]
			if ok {
				for _, report := range reports {
					lat2 := report.location[0]
					long2 := report.location[1]

					distance := QuickDistanceCalc(lat1, long1, lat2, long2)

					(*report).Offset = distance
				}
			}

			reports, ok = reportsByV[node.ID]
			if ok {
				for _, report := range reports {
					lat2 := report.location[0]
					long2 := report.location[1]

					distance := QuickDistanceCalc(lat1, long1, lat2, long2)

					(*report).OffsetRest = distance
				}
			}
		}
		return nil
	}

	log.Infof("updating offsets")
	err := osm.Iterate(OsmNodeIter)
	if err != nil {
		return err
	}
	log.Infof("done updating offsets")

	return nil
}

func downloadOsmPbf() {
	log.Infof("starting osm.pbf download in the background")
	go osm.Download(osm.Config{
		PbfUrl:        config.GetServiceConfig().AggregationConfigs.OsmPbfUrl,
		PbfVersionUrl: config.GetServiceConfig().AggregationConfigs.OsmPbfVersionUrl,
		Username:      config.GetServiceConfig().AggregationConfigs.OsmUser,
		Password:      config.GetServiceConfig().AggregationConfigs.OsmPass,
		Path:          config.GetServiceConfig().AggregationConfigs.OsmDirPath,
		Insecure:      config.GetServiceConfig().AggregationConfigs.InsecureOSMDownload,
	})
}
