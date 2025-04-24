package Redis

import (
    "fmt"
    "github.com/bytedance/sonic"
    log "github.com/sirupsen/logrus"
	"navigation/tools/util"
	"strconv"
    "strings"
)





func convertRedisResultsToFeasibilityConfirmationCases(rawResults map[string]string) (map[string]FeasibilityConfirmationCaseRedisSchema, error) {
	refinedData := make(map[string]FeasibilityConfirmationCaseRedisSchema)
	for id, feasibilityConfirmationCase := range rawResults {
		var singleCase FeasibilityConfirmationCaseRedisSchema
		err := sonic.UnmarshalString(feasibilityConfirmationCase, &singleCase)
		if err != nil {
			return nil, err
		}
		singleCase.Id = id
		refinedData[util.ConcatNumbers(singleCase.U, singleCase.V, singleCase.X)] = singleCase
	}
	return refinedData, nil
}

func convertRedisResultsToAnnotationEvents(rawResults string) AnnotationEvents {
	var annotationEvents AnnotationEvents
	annotationEvents.Events = make(map[string]SegmentsData)
	var results = map[string]string{}
	err := sonic.UnmarshalString(rawResults, &results)
	if err != nil {
		log.WithError(err).Errorf("error converting Redis Result JSON to Traffic : %v", err)
		return annotationEvents
	}
	for segment, events := range results {
		var data SegmentsData
		err := sonic.UnmarshalString(events, &data)
		if err != nil {
			return AnnotationEvents{}
		}
		annotationEvents.Events[segment] = data
	}
	return annotationEvents
}


func getAnnotationEventsKey(keyPattern, version string) string {
	key := fmt.Sprintf("%s:%s", keyPattern, version)
	return key
}


func convertRedisResultsToTraffic(rawResults string) map[TrafficRedisKey]int {
	trafficData := make(map[TrafficRedisKey]int)
	var results = map[string]string{}
	err := sonic.UnmarshalString(rawResults, &results)
	if err != nil {
		log.WithError(err).Errorf("error converting Redis Result JSON to Traffic : %v", err)
		return trafficData
	}
	for segment, segmentCongestion := range results {
		congestion, err := strconv.Atoi(segmentCongestion)
		if err != nil {
			continue
		}

		segmentSplit := strings.Split(segment, "_")

		u, err := strconv.ParseInt(segmentSplit[0], 10, 64)
		if err != nil {
			continue
		}

		v, err := strconv.ParseInt(segmentSplit[1], 10, 64)
		if err != nil {
			continue
		}

		trafficData[TrafficRedisKey{U: u, V: v}] = congestion
	}
	return trafficData
}

func getFeasibilityCasesHashmapKey(baseKey, versionKey string) (string, error) {
	return fmt.Sprintf("%s_v%s",baseKey,versionKey), nil
}

