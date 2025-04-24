package polyline


import (
	"math"
	"strings"
	"navigation/pkg/DTO"
)

func Encode(locations []DTO.IndexedLocation) string {
	var result strings.Builder
	prevLat, prevLng := 0, 0

	for _, loc := range locations {
		lat := int(math.Round(loc.Latitude * 1e5))
		lng := int(math.Round(loc.Longitude * 1e5))

		deltaLat := lat - prevLat
		deltaLng := lng - prevLng

		result.WriteString(encodeSignedInt(deltaLat))
		result.WriteString(encodeSignedInt(deltaLng))

		prevLat = lat
		prevLng = lng
	}

	return result.String()
}

func encodeSignedInt(value int) string {
	shifted := value << 1
	if value < 0 {
		shifted = ^shifted
	}

	var encoded strings.Builder
	for shifted >= 0x20 {
		encoded.WriteByte(byte((shifted&0x1F | 0x20) + 63))
		shifted >>= 5
	}
	encoded.WriteByte(byte(shifted + 63))

	return encoded.String()
}

