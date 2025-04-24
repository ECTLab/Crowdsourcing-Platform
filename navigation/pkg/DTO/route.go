package DTO


type Route struct {
	Geometry   string  `json:"geometry"`
	Legs       []Leg   `json:"legs"`
	WeightName string  `json:"weight_name"`
	Weight     float64 `json:"weight"`
	Duration   float64 `json:"duration"`
	Distance   float64 `json:"distance"`
}

type Leg struct {
	Steps      []Step     `json:"steps"`
	Summary    string     `json:"summary"`
	Weight     float64    `json:"weight"`
	Duration   float64   `json:"duration"`
	Distance   float64   `json:"distance"`
	Annotation Annotation `json:"annotation"`
}

type Step struct {
	Geometry      		string         		`json:"geometry"`
	Maneuver      		Maneuver       		`json:"maneuver"`
	Mode          		string         		`json:"mode"`
	Ref           		string         		`json:"ref"`
	DrivingSide   		string         		`json:"driving_side"`
	Name          		string         		`json:"name"`
	Intersections 		[]Intersection 		`json:"intersections"`
	Weight        		float64        		`json:"weight"`
	Duration      		float64        		`json:"duration"`
	Distance      		float64        		`json:"distance"`
	RotaryName    		string         		`json:"rotary_name"`
	Pronunciation 		string         		`json:"pronunciation"`
	Instruction   		string		 		`json:"instruction"`
	VoiceInstructions 	[]VoiceInstruction 	`json:"voice_instructions"`
}

type VoiceInstruction struct {
	Announcement          string  `json:"announcement"`
	AltAnnouncement       string  `json:"alt_announcement"`
	DistanceAlongGeometry float64 `json:"distance_along_geometry"`
	TimeOffset            float64 `json:"time_offset"`
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
	In		 int	   `json:"in"`
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
	Nodes       				[]int64   `json:"nodes"`
	Distance    				[]float32 `json:"distance"`
	Maxspeed    				[]Maxspeed
	SpeedBump   				[]SpeedBump
	SpeedCamera 				[]SpeedCamera
	Congestion  				[]Congestion
	Police      				[]Police
	FeasibilityConfirmation		[]FeasibilityConfirmation
}

type Maxspeed struct {
	Speed   int32
	Unknown bool
}

type SpeedBump struct {
	Exists bool
	Offset float64
}

type SpeedCamera struct {
	Exists bool
	Offset float64
}

type valueAndText struct {
    Text	string
	Value	float64
}
