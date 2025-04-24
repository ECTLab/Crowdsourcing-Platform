package osrmtool

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
)

func (c *Client) Table(locations []Location, sources, destinations []int) (OsrmTableResponse, error) {
	if len(locations) == 0 {
		return OsrmTableResponse{}, errors.New("called Osrm Table with no locations")
	}

	url := c.url(TableService, DrivingProfile, locations, sources, destinations, false)

	res, err := c.HttpGet(url)
	if err != nil {
		return OsrmTableResponse{}, err
	}

	var response = OsrmTableResponse{}
	err = sonic.Unmarshal(res, &response)
	if err != nil {
		return OsrmTableResponse{}, err
	}

	if response.Code != OkCode {
		return response, errors.New(response.Code + ": " + response.Message)
	}

	return response, nil
}

func (c *Client) Match(locations []Location) (OsrmMatchResponse, error) {
	if len(locations) == 0 {
		return OsrmMatchResponse{}, errors.New("called Osrm Match with no locations")
	}

	url := c.url(MatchService, DrivingProfile, locations, nil, nil, true)

	res, err := c.HttpGet(url)
	if err != nil {
		return OsrmMatchResponse{}, err
	}

	var response = OsrmMatchResponse{}
	err = sonic.Unmarshal(res, &response)
	if err != nil {
		return OsrmMatchResponse{}, err
	}

	if response.Code != OkCode {
		return response, errors.New(response.Code + ": " + response.Message)
	}

	return response, nil
}

func (c *Client) url(service Service, profile Profile, locations []Location, sources, destinations []int, annotations bool) string {
	coords := []string{}
	// TODO add and use other options
	for _, loc := range locations {
		coords = append(coords, fmt.Sprintf("%.8f,%.8f", loc[1], loc[0]))
	}

	coordsJoined := strings.Join(coords, ";")

	base := fmt.Sprintf("%s/%s/v1/%s/%s", c.Address, service, profile, coordsJoined)

	if sources == nil && destinations == nil && !annotations {
		return base
	}

	params := []string{}

	if sources != nil {
		sourcesStr := mapIntToStr(sources)
		param := strings.Join(sourcesStr, ";")
		params = append(params, param)
	}

	if destinations != nil {
		destinationsStr := mapIntToStr(destinations)
		param := strings.Join(destinationsStr, ";")
		params = append(params, param)
	}

	if annotations {
		params = append(params, "annotations=true")
	}

	r := base + "?" + strings.Join(params, "&")

	return r
}

func mapIntToStr(l []int) []string {
	r := []string{}
	for _, n := range l {
		r = append(r, fmt.Sprintf("%d", n))
	}
	return r
}
