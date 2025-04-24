package DTO



type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type IndexedLocation struct {
	Location
	OriginalIndex   int64
}

type Duration struct {
    valueAndText
}

type Distance struct {
    valueAndText
}
type Congestion string

const (
	CongestionUnknown  Congestion = "unknown"
	CongestionLow      Congestion = "low"
	CongestionModerate Congestion = "moderate"
	CongestionHeavy    Congestion = "heavy"
	CongestionSevere   Congestion = "severe"
)

var ToCongestion map[int]Congestion = map[int]Congestion{
	0: CongestionUnknown,
	1: CongestionLow,
	2: CongestionModerate,
	3: CongestionHeavy,
	4: CongestionSevere,
}


type FeasibilityConfirmation struct {
	Id					string	`json:"id"`
	SegmentIndex		int64	`json:"segment_index"`
	Confirmation		bool	`json:"confirmation"`
	Confidence    		float64	`json:"confidence"`
	Contributions 		int		`json:"contributions"`
	RoutingLikelihood 	float64	`json:"routing_likelihood‍"`
}

type Police struct {
	Id           string
	Exists       bool
	Offset       float64
	Confirmation bool
	Confidence   float32
}

type Pronunciation struct {
	Maxspeed float32 `json:"max_speed"`
}


type UserVersion int


const (
	CAB              UserVersion = 2
	CAB_TREATMENT    UserVersion = 5
)


type AnnouncementTypeEnum int

const (
	Unavailable AnnouncementTypeEnum = iota
	Announcement
	AltAnnouncement
)
