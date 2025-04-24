package util

import (
    "bytes"
    "fmt"
	log "github.com/sirupsen/logrus"
	"io"
    "math"
    "navigation/pkg/DTO"
    "os"
    "path/filepath"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"
)



func CreateDirectory() error {
	dirPath := GetCurrentDirectory()
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			log.Println("Error creating directory:", err)
			return err
		}
		log.Println("Directory created successfully")
	} else {
		log.Println("Directory already exists")
	}
	return nil
}

func ConcatNumbers(nodeIds ...int64) string {
	concatinatedString := ""
	for _, node := range nodeIds {
		concatinatedString = concatinatedString + strconv.FormatInt(node, 10) + "."
	}
	return concatinatedString
}


func GetCurrentDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.Join(wd, "/tools/catboost/data")
}

func ConvertLocationToString(location DTO.Location) string {
	return fmt.Sprintf("%f,%f", location.Longitude, location.Latitude)
}

func ConvertLocationsToString(locations []DTO.Location) []string {
	var result []string
	for _, loc := range locations {
		result = append(result, fmt.Sprintf("%f,%f", loc.Longitude, loc.Latitude))
	}
	return result
}

func InverseTransform(firstBlockEta, targetTransformed float64) float32 {
	return float32(firstBlockEta * targetTransformed)
}

func GetTimeOfDay(timeOfDay time.Time) float32 {
	return float32(timeOfDay.Hour()) + float32(math.Floor(float64(timeOfDay.Minute())/15)/4)
}


type Polyline struct {
	Points string `json:"points"`
}

// DecodePolyline converts a polyline encoded string to an array of LatLng objects.
func DecodePolyline(poly string) ([]LatLng, error) {
	p := &Polyline{
		Points: poly,
	}
	return p.Decode()
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Decode converts this encoded Polyline to an array of LatLng objects.
func (p *Polyline) Decode() ([]LatLng, error) {
	input := bytes.NewBufferString(p.Points)

	var lat, lng int64
	path := make([]LatLng, 0, len(p.Points)/2)
	for {
		dlat, _ := decodeInt(input)
		dlng, err := decodeInt(input)
		if err == io.EOF {
			return path, nil
		}
		if err != nil {
			return nil, err
		}

		lat, lng = lat+dlat, lng+dlng
		path = append(path, LatLng{
			Lat: float64(lat) * 1e-5,
			Lng: float64(lng) * 1e-5,
		})
	}
}

// decodeInt reads an encoded int64 from the passed io.ByteReader.
func decodeInt(r io.ByteReader) (int64, error) {
	result := int64(0)
	var shift uint8

	for {
		raw, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		b := raw - 63
		result += int64(b&0x1f) << shift
		shift += 5

		if b < 0x20 {
			bit := result & 1
			result >>= 1
			if bit != 0 {
				result = ^result
			}
			return result, nil
		}
	}
}


func GetLocalTime(t time.Time) time.Time {
	utc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		log.Error(err)
		return time.Now()
	}
	return t.In(utc)
}


func CompareStringNumbers(str1, str2 string) (int, error) {
	num1, err := strconv.Atoi(str1)
	if err != nil {
		return 0, fmt.Errorf("first string is not a valid integer: %v", err)
	}
	num2, err := strconv.Atoi(str2)
	if err != nil {
		return 0, fmt.Errorf("second string is not a valid integer: %v", err)
	}
	if num1 < num2 {
		return -1, nil
	} else if num1 > num2 {
		return 1, nil
	} else {
		return 0, nil
	}
}

func IsTheBearingValid(androidVersion string) bool {
	versionNumbers := strings.Split(androidVersion, ".")
	eligibleAndroidVersionNumbers := strings.Split("7.7.0", ".")
	if len(eligibleAndroidVersionNumbers) != 3 || len(versionNumbers) != 3 {
		return false
	}
	for i := 0 ; i < len(versionNumbers) ; i++ {
		comparisonResult, err := CompareStringNumbers(versionNumbers[i], eligibleAndroidVersionNumbers[i])
		if err != nil {
			return false
		}
		if comparisonResult == -1 {
			return false
		} else if comparisonResult == 1 {
			return true
		}
	}
	return true
}

func CountElementsInBoolArray(arr []bool, target bool) int {
	cnt := 0
	for _, element := range arr {
		if element == target {
			cnt++
		}
	}
	return cnt
}

func AngularDistance(angle1, angle2 int) float64 {
	diff := int(math.Abs(float64(angle1-angle2))) % 360
    return math.Min(float64(diff), float64(360-diff))
}