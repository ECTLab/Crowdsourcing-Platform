package Route

import (
	"navigation/pkg/DTO"
)


type StreamNavigationResponse struct {
	RequestId string        `json:"request_id"`
	Routes    []StreamRoute `json:"routes"`
	TTL       int64         `json:"ttl"`
}

type StreamRoute struct {
	Geometry   string      `json:"geometry"`
	Legs       []StreamLeg `json:"legs"`
	WeightName string      `json:"weight_name"`
	Weight     float64 			`json:"weight"`
	Duration   float64 			`json:"duration"`
	Distance   float64 			`json:"distance"`
}

type StreamLeg struct {
	Steps      []StreamStep `json:"steps"`
	Summary    string       `json:"summary"`
	Weight     float64    		`json:"weight"`
	Duration   float64   		`json:"duration"`
	Distance   float64          `json:"distance"`
	Annotation StreamAnnotation `json:"annotation"`
}

type StreamStep struct {
	Geometry      		string  	       		`json:"geometry"`
	Maneuver      		DTO.Maneuver       		`json:"maneuver"`
	Name          		string      	   		`json:"name"`
	Intersections 		[]DTO.Intersection 		`json:"intersections"`
	Weight        		float64        			`json:"weight"`
	Duration      		float64        			`json:"duration"`
	Distance      		float64        			`json:"distance"`
	RotaryName    		string         			`json:"rotary_name"`
	VoiceInstructions 	[]DTO.VoiceInstruction 	`json:"voice_instructions"`
}

type StreamAnnotation struct {
	PoliceSegmentIndices		[]int							`json:"police_segment_indices"`
	SpeedBumpSegmentIndices     []int							`json:"speed_bump_segment_indices"`
	SpeedCameraSegmentIndices   []int							`json:"speed_camera_segment_indices"`
	FeasibilityConfirmation 	[]DTO.FeasibilityConfirmation	`json:"feasibility_confirmation"`
}