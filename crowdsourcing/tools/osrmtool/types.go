package osrmtool

type Profile string

const (
	DrivingProfile = "driving"
)

type Service string

const (
	RouteService   = "route"
	NearestService = "nearest"
	TableService   = "table"
	MatchService   = "match"
)

type Location [5]float64 // [lat, lng, bearing, timestamp, accuracy]

// OSRM types

const OkCode string = "Ok"

type OsrmResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type OsrmTableResponse struct {
	OsrmResponse
	Durations [][]float64 `json:"durations"`
	// Sources      []Waypoint `json:"sources"`
	// Destinations []Waypoint `json:"destinations"`
}

type OsrmMatchResponse struct {
	OsrmResponse
	Tracepoints []Waypoint `json:"tracepoints"`
	Matchings   []Route    `json:"matchings"`
}

type Route struct {
	Geometry   string  `json:"geometry"`
	Legs       []Leg   `json:"legs"`
	WeightName string  `json:"weight_name"`
	Weight     float64 `json:"weight"`
	Duration   float64 `json:"duration"`
	Distance   float64 `json:"distance"`
	Confidence float64 `json:"confidence,omitempty"`
}

type Leg struct {
	Steps      []Step     `json:"steps"`
	Summary    string     `json:"summary"`
	Weight     float64    `json:"weight"`
	Duration   float64    `json:"duration"`
	Distance   float64    `json:"distance"`
	Annotation Annotation `json:"annotation"`
}

type Step struct {
	Geometry      string         `json:"geometry"`
	Maneuver      Maneuver       `json:"maneuver"`
	Mode          string         `json:"mode"`
	Ref           string         `json:"ref"`
	DrivingSide   string         `json:"driving_side"`
	Name          string         `json:"name"`
	Intersections []Intersection `json:"intersections"`
	Weight        float64        `json:"weight"`
	Duration      float64        `json:"duration"`
	Distance      float64        `json:"distance"`
	RotaryName    string         `json:"rotary_name"`
	Pronunciation string         `json:"pronunciation"`
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
	Hint     string    `json:"hint"`
	Distance float64   `json:"distance"`
	Name     string    `json:"name"`
	Location []float64 `json:"location"`
}

type Annotation struct {
	Nodes    []int64   `json:"nodes"`
	Distance []float32 `json:"distance"`
}
