package Osrm

import (
    "encoding/json"
    "fmt"
    "math"
    "math/rand"
    httpClient "navigation/pkg/Clients/Http"
    "navigation/pkg/DTO"
    "net/http"
    "time"
)


type tripResponse struct {
	Code      string     `json:"code"`
	Trips     []trip     `json:"trips"`
	Waypoints []Waypoint `json:"waypoints"`
}

type trip struct {
	Geometry   string 		`json:"geometry"`
	Legs       []DTO.Leg  	`json:"legs"`
	WeightName string 		`json:"weight_name"`
	Weight     float64 		`json:"weight"`
	Duration   float64 		`json:"duration"`
	Distance   float64 		`json:"distance"`
}

type Client struct {
	Version        string
	Profile        string
	Timeout        time.Duration
	MaxRetries     int
	BackoffMax     time.Duration
	BackoffFactor  float64
	BackoffEnabled bool
	Session        *http.Client
	DecodeOutput   bool
}

type BaseRequest struct {
	Coordinates []DTO.Location
	Radiuses    []float64
	Bearings    []Bearing
	Hints       []string
	Service     string
}

type RouteRequest struct {
	*BaseRequest
	Origin         DTO.Location
	Destination    DTO.Location
	Waypoints      []DTO.Location
	ExclusionZones []string
	Alternative    bool
	Bearing        int
	IsReRoute      bool
}

type BaseResponse struct {
	Code    string
	Message string
}

type TableResponse struct {
	*BaseResponse
	Durations    []float64
	Distances    []float64
	Sources      []*Waypoint
	Destinations []*Waypoint
}

type NavigationService struct {
	Name                            string
	Host                 string
	Client               Client
	DefaultServiceClient *http.Client
}


type NavigationRequest struct {
	Origin         DTO.Location
	Destination    DTO.Location
	Waypoints      []DTO.Location
	ExclusionZones []string
	Alternative    bool
	Bearing        int64
	IsReRoute      bool
}

type NearestStreetRequest struct {
	Location DTO.Location
	Number   int64
}

type NearestStreetResponse struct {
	Code      string     `json:"code"`
	Waypoints []Waypoint `json:"waypoints"`
}

type MatchResponse struct {
	Code        string        `json:"code"`
	Tracepoints []Tracepoint  `json:"tracepoints"`
	Matchings   []Matching    `json:"matchings"`
}

type Tracepoint struct {
	AlternativesCount int       `json:"alternatives_count"`
	WaypointIndex     int       `json:"waypoint_index"`
	MatchingsIndex    int       `json:"matchings_index"`
	Location          []float64 `json:"location"`
	Name              string    `json:"name"`
	Hint              string    `json:"hint"`
	Distance          float64   `json:"distance"`
}

type Matching struct {
	Confidence float64     	`json:"confidence"`
	Geometry   string   	`json:"geometry"`
	Distance   float64  	`json:"distance"`
	Duration   float64  	`json:"duration"`
	Weight     float64  	`json:"weight"`
	WeightName string       `json:"weight_name"`
	Legs       []DTO.Leg    `json:"legs"`
}
type NavigationResponse struct {
	Code      string      `json:"code"`
	Routes    []DTO.Route `json:"routes"`
	Waypoints []Waypoint  `json:"waypoints"`
}

type Maneuver struct {
	Exit          int       `json:"exit"`
	BearingAfter  int       `json:"bearing_after"`
	BearingBefore int       `json:"bearing_before"`
	Location      []float32 `json:"location"`
	Modifier      string    `json:"modifier"`
	Type          string    `json:"type"`
}

type Intersection struct {
	Out      int       `json:"out"`
	Entry    []bool    `json:"entry"`
	Bearings []int     `json:"bearings"`
	Location []float64 `json:"location"`
}

type Waypoint struct {
	WaypointIndex int       `json:"waypoint_index,omitempty"`
	TripsIndex    int       `json:"trips_index,omitempty"`
	Hint     string    `json:"hint"`
	Distance float64   `json:"distance"`
	Name     string    `json:"name"`
	Location []float64 `json:"location"`
}

type Annotation struct {
	Nodes       []int64   `json:"nodes"`
	Distance    []float32 `json:"distance"`
	Maxspeed    []DTO.Maxspeed
	SpeedBump   []DTO.SpeedBump
	SpeedCamera []DTO.SpeedCamera
	Congestion  []DTO.Congestion
	Police      []DTO.Police
}

type Bearing struct {
	value  int
	range_ int
}

type Pronunciation struct {
	Maxspeed float32 `json:"max_speed"`
}

func (c Client) expBackoff(attempt int) time.Duration {
	timeout := c.Timeout * time.Duration(math.Pow(2, float64(attempt)))
	if timeout > c.BackoffMax {
		timeout = c.BackoffMax
	}
	jitter := time.Duration(c.BackoffFactor * float64(rand.Int63n(int64(c.Timeout))))
	return timeout + jitter
}

func (c Client) request(url string) (interface{}, error) {
	attempt := 0
	for attempt < c.MaxRetries {
		body, err := httpClient.Get(url, c.Session)
		if err != nil {
			if !c.BackoffEnabled {
				return nil, err
			}
			time.Sleep(c.expBackoff(attempt))
			attempt++
			continue
		}
		if c.DecodeOutput {
			var decodedResponse TableResponse
			if err := json.Unmarshal(body, &decodedResponse); err != nil {
				return nil, err
			}
			return &decodedResponse, nil
		} else {
			var result map[string]interface{}
			if err := json.Unmarshal(body, &result); err != nil {
				return nil, err
			}
			return result, nil
		}
	}
	return nil, fmt.Errorf("max retries reached")
}

func (c Client) OsrmCall(url string) (interface{}, error) {
	return c.request(url)
}

func (c Client) Close() {
	c.Session.CloseIdleConnections()
}


func NewOSRMRouteRequest(origin DTO.Location, destination DTO.Location, waypoints []DTO.Location) * RouteRequest {
	return &RouteRequest{
		Origin:         origin,
		Destination:    destination,
		Waypoints:      waypoints,
	}
}


