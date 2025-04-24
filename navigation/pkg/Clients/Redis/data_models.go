package Redis

import "time"

type AnnotationEvents struct {
	Events map[string]SegmentsData
}

type SegmentsData struct {
	Length float64      `json:"length"`
	Events []EventsData `json:"events"`
}

type EventsData struct {
	Event  string  `json:"event"`
	OffSet float64 `json:"offset"`
}

var AnnotationEventsData AnnotationEvents

var CurrentVersion time.Time
var SpeedBumpName = "speed_bump"
var SpeedCameraName = "speed_camera"

type TrafficRedisKey struct {
	U, V int64
}

var TrafficData map[TrafficRedisKey]int

// version is like "annotation_traffic_segments:2024-03-16-12-45"
var CurrentTrafficVersion string = ""

type RedisKey struct {
	U, V int64
}

// all police keys are stored in order.
// U is always greater than V
type PoliceRedisKey RedisKey
type AccidentRedisKey RedisKey

type OnlineReportType string

const (
	POLICE_ONLINE   			OnlineReportType = "POLICE"
	ACCIDENT_ONLINE 			OnlineReportType = "ACCIDENT"
	FEASIBILITY_CONFIRMATION 	OnlineReportType = "FEASIBILITY"
)

type OnlineReportRedisSchema struct {
	Confirmation bool             `json:"confirmation,omitempty"`
	Confidence   float32          `json:"confidence,omitempty"`
	Type         OnlineReportType `json:"type,omitempty"`
	U            int64            `json:"u,omitempty"`
	V            int64            `json:"v,omitempty"`
	Offset       float64          `json:"offset,omitempty"`
	OffsetRest   float64          `json:"offset_rest,omitempty"`
	Uid          string           `json:"uid,omitempty"`
}

type FeasibilityConfirmationCaseRedisSchema struct {
	Id				  string            `json:"id,omitempty"`
	Confirmation      bool				`json:"confirmation,omitempty"`
  	Confidence        float64			`json:"confidence,omitempty"`
  	Type              OnlineReportType	`json:"type,omitempty"`
  	U                 int64				`json:"u,omitempty"`
  	V                 int64				`json:"v,omitempty"`
  	X                 int64				`json:"x,omitempty"`
  	Contributions     int				`json:"contributions,omitempty"`
  	RoutingLikelihood float64			`json:"routing_likelihood,omitempty"`
}

var PoliceData map[PoliceRedisKey]OnlineReportRedisSchema
var AccidentData map[AccidentRedisKey]OnlineReportRedisSchema
var FeasibilityConfirmationData map[string]FeasibilityConfirmationCaseRedisSchema

var CurrentOnlineReportsVersion string = ""
