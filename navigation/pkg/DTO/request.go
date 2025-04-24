package DTO


type NavigationRequest struct {
	VehicleType        string	`json:"vehicleType"`
	Origin             Location	`json:"origin"`
	Destination        Location `json:"destination"`
	MiddleDestinations []Location `json:"middleDestinations"`
	RouteReason        string `json:"routeReason"`
}

type Profile string

const (
	ProfileUnspecified = "UNSPECIFIED"
	ProfileCar = "CAR"
	Profilebike = "BIKE"
)

type RouteRequestReason int

const (
	RouteRequestReasonUnknown RouteRequestReason = iota
	RouteRequestReasonRideStart
	RouteRequestReasonOutOfRoute
	RouteRequestReasonTtlExpired
	RouteRequestReasonSettingChanged
)

var reasonsConsideredAsReroute = []RouteRequestReason{
	RouteRequestReasonOutOfRoute,
	RouteRequestReasonTtlExpired,
	RouteRequestReasonSettingChanged,
 }

func (rrr RouteRequestReason) Name() string {
	switch rrr {
	case RouteRequestReasonRideStart:
		return "RideStart"
	case RouteRequestReasonOutOfRoute:
		return "OutOfRoute"
	case RouteRequestReasonTtlExpired:
		return "TtlExpired"
	case RouteRequestReasonSettingChanged:
		return "SettingChanged"
	default:
		return "Unknown"
	}
}


type TTSRequest struct {
	RequestId        string `json:"request_id"`
	Text             string `json:"text"`
	Tid              string `json:"tid"`
	UseCDN           bool   `json:"use_cdn"`
	ExpirationPolicy string `json:"expiration_policy"`
}


type GetVoiceLinkRequest struct {
	Announcement    string `json:"announcement"`
	AltAnnouncement string `json:"alt_announcement"`
}

type GetNearestStreetLocationRequest struct {
	Location 		Location	`json:"location"`
	Number			int64		`json:"number"`
	ZoomLevel		float32		`json:"zoomLevel"`
	LimitCities     bool        `json:"limitCities"`
}

type ConvertTextToSpeechRequest struct {
    ClientID        string `json:"client_id"`
    Text            string `json:"text"`
    ExpirationPolicy string `json:"expiration_policy"`
}

type MatchingRequest struct {
	Locations		[]Location `json:"locations"`
}

type MatchingResponse struct {
	MatchedLocations  []IndexedLocation `json:"matched_locations"`
	GeometryPolyline  string            `json:"geometry_polyline"`
}

type TspRequest struct {
	Locations		[]Location 	`json:"locations"`
	RoundTrip 		bool		`json:"round_trip"`
	StartFromFirst  bool		`json:"start_from_first"`
	FinishInLast	bool		`json:"finish_in_last"`
	Profile			Profile		`json:"profile"`
}

type TspResponse struct {
	Points			[]IndexedLocation	`json:"points"`
	Duration		int64				`json:"duration"`
	Distance		int64				`json:"distance"`
}