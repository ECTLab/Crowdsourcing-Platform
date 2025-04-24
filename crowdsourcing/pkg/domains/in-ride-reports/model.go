package inridereports


type InRideReportType string

const (
	UNSPECIFIED InRideReportType = "UNSPECIFIED"
	DEAD_END   	InRideReportType = "DEAD_END"
	NOT_EXISTS 	InRideReportType = "NOT_EXISTS"
	OFF_ROAD   	InRideReportType = "OFF_ROAD"
	WRONG_TURN 	InRideReportType = "WRONG_TURN"
	NO_ENTRY   	InRideReportType = "NO_ENTRY"
	NO_CAR     	InRideReportType = "NO_CAR"
	POLICE     	InRideReportType = "POLICE"
	ACCIDENT   	InRideReportType = "ACCIDENT"
	FEASIBILITY	InRideReportType = "FEASIBILITY"
	TEMPORARILY_CLOSED InRideReportType = "TEMPORARILY_CLOSED"
	TRAFFIC            InRideReportType = "TRAFFIC"
	CAMERA             InRideReportType = "CAMERA"
	SPEED_BUMP         InRideReportType = "SPEED_BUMP"
)

type OnlineReportType string

const (
	POLICE_ONLINE   OnlineReportType = OnlineReportType(POLICE)
	ACCIDENT_ONLINE OnlineReportType = OnlineReportType(ACCIDENT)
)


var OnlineReportTypes = []OnlineReportType{
	POLICE_ONLINE,
	ACCIDENT_ONLINE,
}

type LocationTime struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Time      int64   `json:"time,omitempty"`
}

type Reporter struct {
	Id     		string `json:"id"`
	DeviceId 	string `json:"device_id"`
	Role   		string `json:"role,omitempty"`
	Client 		string `json:"client,omitempty"`
}

type InRideReport struct {
	Type                   InRideReportType `json:"type,omitempty"`
	EngagementLocationTime LocationTime     `json:"engagement_location_time,omitempty"`
	SubmitLocationTime     LocationTime     `json:"submit_location_time,omitempty"`
	Geometry               string           `json:"geometry,omitempty"`
	Footprint              string           `json:"footprint"`
	Online                 bool             `json:"online,omitempty"`
	RideId                 string           `json:"ride_id"`
	NavigationRequestId    string           `json:"navigation_request_id,omitempty"`
	Reporter               Reporter         `json:"reporter,omitempty"`
	TripId				   string			`json:"trip_id"`
}

type InRideReportResponse struct {
	ReportId string `json:"report_id"`
}

type InRideKafkaReport struct {
	InRideReport
	CreatedAt string `json:"CreatedAt"`
}

type RedisLocation [2]float64 // [Latitude, Longitude]

type Event struct {
	PositiveConfirmation int32            `json:"positive_confirmation,omitempty"`
	NegativeConfirmation int32            `json:"negative_confirmation,omitempty"`
	Type                 OnlineReportType `json:"type,omitempty"`
	Location             RedisLocation    `json:"location,omitempty"`
	U                    int64            `json:"u,omitempty"`
	V                    int64            `json:"v,omitempty"`
	Offset               float64          `json:"offset,omitempty"`
	OffsetRest           float64          `json:"offset_rest,omitempty"`
	H3Index              string           `json:"h3_index,omitempty"`
	Uid                  string           `json:"uid,omitempty"`
	ReportTimestamp      int64            `json:"report_time,omitempty"`
	UpdateTimestamp      int64            `json:"update_time,omitempty"`
	TtlMinutes           int32            `json:"ttl,omitempty"`
	IsAggregated         bool             `json:"is_aggregated,omitempty"`
}

type ReportConfirmation struct {
	ReportId  string   			`json:"id,omitempty"`
	Confirmed bool     			`json:"confirmed,omitempty"`
	UserResponse string         `json:"user_response,omitempty"`
	Type      InRideReportType 	`json:"type,omitempty"`
	Reporter  Reporter 			`json:"reporter,omitempty"`
}

type ReportConfirmationResponse struct {

}

type AnnotationOnlineReport struct {
	Confidence   float64          `json:"confidence,omitempty"`
	Confirmation bool             `json:"confirmation,omitempty"`
	Type         OnlineReportType `json:"type,omitempty"`
	U            int64            `json:"u,omitempty"`
	V            int64            `json:"v,omitempty"`
	Offset       float64          `json:"offset"`
	OffsetRest   float64          `json:"offset_rest"`
	Uid          string           `json:"uid,omitempty"`
	location     RedisLocation
}

type AnnotationFeasibilityConfirmationCase struct {
	Confirmation      bool				`json:"confirmation,omitempty"`
  	Confidence        float64			`json:"confidence,omitempty"`
  	Type              InRideReportType	`json:"type,omitempty"`
  	U                 int64				`json:"u,omitempty"`
  	V                 int64				`json:"v,omitempty"`
  	X                 int64				`json:"x,omitempty"`
  	Contributions     int				`json:"contributions,omitempty"`
  	RoutingLikelihood float64			`json:"routing_likelihood,omitempty"`
}

type FeasibilityConfirmationKafka struct {
	Type             InRideReportType    	`json:"type,omitempty"`
  	U                int64					`json:"u,omitempty"`
  	V                int64					`json:"v,omitempty"`
  	X                int64					`json:"x,omitempty"`
  	Contributions    int					`json:"contributions,omitempty"`
  	Id               string					`json:"id,omitempty"`
  	Confirmed        bool					`json:"confirmed,omitempty"`
	UserResponse     string                  `json:"userResponse,omitempty"`
	CreatedAt        string					`json:"createdAt,omitempty"`
}
