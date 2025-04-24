package h3

import (
	"errors"

	"github.com/uber/h3-go/v4"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var Resolution int = 9

func GetH3CellWithResolution(lat float64, lng float64, resolution int) h3.Cell {
	latLng := h3.NewLatLng(lat, lng)
	return h3.LatLngToCell(latLng, resolution)
}

func GetH3CellFromStringIndex(str string) (h3.Cell, error) {
	cell := h3.Cell(h3.IndexFromString(str))
	ok := cell.IsValid()
	if !ok {
		return 0, errors.New("h3: provided index (" + str + ") is not a valid H3 index")
	}
	return cell, nil
}

func WithNeighbours(cell h3.Cell) []h3.Cell {
	return cell.GridDisk(1)
}

func AsIndex(cells []h3.Cell) []string {
	var indexes = []string{}
	for _, cell := range cells {
		indexes = append(indexes, cell.String())
	}
	return indexes
}

func GetH3Index(lat float64, lng float64) string {
	return GetH3CellWithResolution(lat, lng, Resolution).String()
}

func GetH3IndexWithNeighbours(lat float64, lng float64) [7]string {
	cell := GetH3CellWithResolution(lat, lng, Resolution)
	cells := WithNeighbours(cell)
	indexSlice := AsIndex(cells)

	var indexes = [7]string{}

	if len(indexSlice) == 7 {
		for i := 0; i < 7; i++ {
			indexes[i] = cells[i].String()
		}
	}

	return indexes
}
