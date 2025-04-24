package Redis

import (
    "context"
    "errors"
    "github.com/bytedance/sonic"
    "github.com/go-redis/redis/v8"
    log "github.com/sirupsen/logrus"
)

func (c *client) GetLatestAnnotationEvents(ctx context.Context, version string) (AnnotationEvents, error) {
	results, err := c.RedisClient.Get(ctx, getAnnotationEventsKey("annotation_events_segments", version)).Result()

	annotationEvents := convertRedisResultsToAnnotationEvents(results)

	if errors.Is(err, redis.Nil) {
		return AnnotationEvents{}, redis.Nil
	}
	if err != nil {
		return AnnotationEvents{}, err
	}
	return annotationEvents, nil
}

func (c *client) GetAnnotationLatestVersion(ctx context.Context) (string, error) {
	latestVersion, err := c.RedisClient.Get(ctx, "annotation_events_latest_version").Result()
	if errors.Is(err, redis.Nil) {
		return "", redis.Nil
	}
	if err != nil {
		return "", err
	}
	return latestVersion, nil
}


func (c *client) GetLatestTrafficData(ctx context.Context, trafficLatestVersion string) (map[TrafficRedisKey]int, error) {

	key := trafficLatestVersion
	results, err := c.RedisClient.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, redis.Nil
	}
	if err != nil {
		return nil, err
	}

	trafficData := convertRedisResultsToTraffic(results)

	return trafficData, nil
}

func (c *client) GetTrafficLatestVersion(ctx context.Context) (string, error) {
	version, err := c.RedisClient.Get(ctx, "annotation_traffic_latest_version").Result()
	if errors.Is(err, redis.Nil) {
		return "", redis.Nil
	}
	if err != nil {
		return "", err
	}
	return version, nil
}


func (c *client) GetOnlineReportsLatestVersion(ctx context.Context) (string, error) {
	versionRedisKey := "annotation_online_reports_latest_version"
	version, err := c.RedisClient.Get(ctx, versionRedisKey).Result()
	if errors.Is(err, redis.Nil) {
		return "", redis.Nil
	}
	if err != nil {
		return "", err
	}
	return version, nil
}

func (c *client) GetLatestOnlineReportsData(ctx context.Context) (
	map[PoliceRedisKey]OnlineReportRedisSchema,
	map[AccidentRedisKey]OnlineReportRedisSchema,
	error,
) {

	onlineReportRedisKey := "online_reports"

	redisResults, err := c.RedisClient.Get(ctx, onlineReportRedisKey).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil, redis.Nil
	}
	if err != nil {
		return nil, nil, err
	}

	var policeData = map[PoliceRedisKey]OnlineReportRedisSchema{}
	var accidentData = map[AccidentRedisKey]OnlineReportRedisSchema{}

	var onlineReportData = []OnlineReportRedisSchema{}
	err = sonic.UnmarshalString(redisResults, &onlineReportData)
	if err != nil {
		log.WithError(err).Errorf("error converting Redis Result JSON to Online Report Type : %v", err)
		return policeData, accidentData, err
	}

	for _, report := range onlineReportData {
		if report.Type == POLICE_ONLINE {
			key := PoliceRedisKey{U: report.U, V: report.V}
			policeData[key] = report
		}
		if report.Type == ACCIDENT_ONLINE {
			key := AccidentRedisKey{U: report.U, V: report.V}
			accidentData[key] = report
		}
	}

	return policeData, accidentData, nil
}

func (c *client) GetLatestFeasibilityConfirmationCases(ctx context.Context) ( map[string]FeasibilityConfirmationCaseRedisSchema, error ) {


	redisFeasibilityCasesHashmapKeyBase := "feasibility_confirmation_cases"
	redisFeasibilityCasesLatestVersionKey := "feasibility_confirmation_latest_version"

	version, err := c.RedisClient.Get(ctx, redisFeasibilityCasesLatestVersionKey).Result()
	if err != nil {
		return nil, err
	}
	key, err := getFeasibilityCasesHashmapKey(redisFeasibilityCasesHashmapKeyBase, version)
	redisResults, err := c.RedisClient.HGetAll(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, redis.Nil
	}
	if err != nil {
		return nil, err
	}


	feasibilityConfirmationCases, err := convertRedisResultsToFeasibilityConfirmationCases(redisResults)

	if errors.Is(err, redis.Nil) {
		return nil, redis.Nil
	}
	if err != nil {
		return nil, err
	}
	return feasibilityConfirmationCases, nil
}