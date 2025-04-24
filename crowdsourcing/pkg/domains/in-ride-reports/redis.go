package inridereports

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"crowdsourcing/config"
	"crowdsourcing/pkg/clients/redis"

	"github.com/bytedance/sonic"
)


const redisEventKeyTemplate string = "{h3}:{event_type}:{uid}"

type redisEventKey struct {
	h3        string
	eventType OnlineReportType
	uid       string
}

// implements fmt.Stringer
func (k redisEventKey) String() string {
	str := strings.ReplaceAll(redisEventKeyTemplate, "{h3}", k.h3)
	str = strings.ReplaceAll(str, "{event_type}", string(k.eventType))
	str = strings.ReplaceAll(str, "{uid}", k.uid)
	return str
}

func (k *redisEventKey) From(str string) error {
	parts := strings.Split(str, ":")

	if len(parts) != 3 {
		return errors.New("redisEventKey: Could not parse from string(" + str + ")")
	}

	k.h3 = parts[0]
	k.eventType = OnlineReportType(parts[1])
	k.uid = parts[2]

	return nil
}

const redisEventHashName string = "events_hash"
const redisNewEventListName string = "new_events_list"
const redisAnnotationReports string = "online_reports"
const redisFeasibilityCasesHashmapKeyBase = "feasibility_confirmation_cases"
const redisAnnotationReportsVersion string = "annotation_online_reports_latest_version"
const redisFeasibilityCasesLatestVersionKey = "feasibility_confirmation_latest_version"

func RedisAddNewOnlineReport(event Event) error {
	event_str, err := sonic.MarshalString(&event)
	if err != nil {
		return err
	}

	eventKey := redisEventKey{
		h3:        event.H3Index,
		eventType: event.Type,
		uid:       event.Uid,
	}
	redisKey := eventKey.String()

	err = redis.Crowdsourcing.ListPush(redisNewEventListName, redisKey)
	if err != nil {
		return err
	}

	err = redis.Crowdsourcing.HSet(redisEventHashName, redisKey, event_str)
	if err != nil {
		return err
	}

	return nil
}

func RedisReadAllOnlineReports() (map[OnlineReportType]map[string][]Event, error) {
	events := map[OnlineReportType]map[string][]Event{
		POLICE_ONLINE:   {},
		ACCIDENT_ONLINE: {},
	}

	rawEvents, err := redis.Crowdsourcing.HGetAll(redisEventHashName)
	if err != nil {
		return nil, err
	}

	for key, event_str := range rawEvents {
		var event Event
		var redisKey redisEventKey
		err := redisKey.From(key)
		if err != nil {
			return nil, err
		}

		err = sonic.UnmarshalString(event_str, &event)
		if err != nil {
			return nil, err
		}

		if event.H3Index != redisKey.h3 || event.Type != redisKey.eventType || event.Uid != redisKey.uid {
			return nil, errors.New("Online Report redis key and report fields mismatch: " + key + ", " + event_str)
		}

		events[redisKey.eventType][redisKey.h3] = append(events[redisKey.eventType][redisKey.h3], event)
	}

	return events, nil
}

func RedisReadAllNewOnlineReports() (map[OnlineReportType]map[string][]Event, []redisEventKey, error) {
	keys, err := redis.Crowdsourcing.ListRangeAll(redisNewEventListName)
	if err != nil {
		return nil, nil, err
	}

	err = redis.Crowdsourcing.ListLPopN(redisNewEventListName, len(keys))
	if err != nil {
		return nil, nil, err
	}

	events, err := RedisReadAllOnlineReports()
	if err != nil {
		return nil, nil, err
	}

	newEventKeys := []redisEventKey{}
	for _, key := range keys {
		var redisKey redisEventKey
		err := redisKey.From(key)
		if err != nil {
			return nil, nil, err
		}

		newEventKeys = append(newEventKeys, redisKey)
	}

	return events, newEventKeys, nil
}

func RedisReplaceAllOnlineReports(reports map[OnlineReportType]map[string][]Event) error {
	err := redis.Crowdsourcing.Del(redisEventHashName)
	if err != nil {
		return err
	}

	for reportType, reportsByType := range reports {
		for h3index, cell := range reportsByType {
			for _, report := range cell {

				report_str, err := sonic.MarshalString(&report)
				if err != nil {
					return err
				}

				eventKey := redisEventKey{
					h3:        h3index,
					eventType: reportType,
					uid:       report.Uid,
				}
				redisKey := eventKey.String()

				err = redis.Crowdsourcing.HSet(redisEventHashName, redisKey, report_str)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func RedisAnnotationUpdateReports(events []*AnnotationOnlineReport) error {
	ttl := config.GetServiceConfig().RedisAnnotation.EventDefaultTtlMinutes

	events_str, err := sonic.MarshalString(&events)
	if err != nil {
		return err
	}

	err = redis.Annotation.SetWithTTL(redisAnnotationReports, events_str, time.Minute*time.Duration(ttl))
	if err != nil {
		return err
	}

	err = redis.Annotation.Set(redisAnnotationReportsVersion, fmt.Sprintf("%d", time.Now().UnixMilli()))
	if err != nil {
		return err
	}

	return nil
}

var ErrNotExists = errors.New("not exists")

func RedisGetOnlineReport(reportId string) (Event, error) {
	var event Event

	event_str, err := redis.Crowdsourcing.HGet(redisEventHashName, reportId)
	if err != nil {
		return event, err
	}

	if event_str == "" {
		return event, ErrNotExists
	}

	err = sonic.UnmarshalString(event_str, &event)
	if err != nil {
		return event, err
	}

	return event, nil
}

func RedisSetOnlineReport(reportId string, event Event) error {
	event_str, err := sonic.MarshalString(&event)
	if err != nil {
		return err
	}

	err = redis.Crowdsourcing.HSet(redisEventHashName, reportId, event_str)
	if err != nil {
		return err
	}

	return nil
}

func RedisGetFeasibilityConfirmationCase(reportId string) (AnnotationFeasibilityConfirmationCase, error) {
	hashmapKey, err := getFeasibilityCasesHashmapKey()
	if err != nil {
		return AnnotationFeasibilityConfirmationCase{}, err
	}
	var event AnnotationFeasibilityConfirmationCase
	event_str, err := redis.Annotation.HGet(hashmapKey, reportId)
	if err != nil {
		return event, err
	}

	err = sonic.UnmarshalString(event_str, &event)
	if err != nil {
		return event, err
	}

	return event, nil
}

func RedisSetFeasibilityConfirmationCase(reportId string, event AnnotationFeasibilityConfirmationCase) error {
	hashmapKey, err := getFeasibilityCasesHashmapKey()
	if err != nil {
		return err
	}
	eventStr, err := sonic.MarshalString(&event)
	if err != nil {
		return err
	}
	err = redis.Annotation.HSet(hashmapKey, reportId, eventStr)
	if err != nil {
		return err
	}

	return nil
}


func getFeasibilityCasesHashmapKey() (string, error) {
	version, err := redis.Annotation.Get(redisFeasibilityCasesLatestVersionKey)
	if err != nil {
		return "", errors.New("faced erro while trying to get latest feasibility versions: " + err.Error())
	}
	return fmt.Sprintf("%s_v%s",redisFeasibilityCasesHashmapKeyBase,version), nil
}