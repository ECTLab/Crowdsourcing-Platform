package Osrm

import (
    "encoding/json"
    "navigation/pkg/Clients/Http"
    "navigation/pkg/DTO"

    log "github.com/sirupsen/logrus"
)

var NavigationOsrmService NavigationService
var navigationOsrmNewClient Client
func (ons NavigationService) GetProviderName() string {
	return ons.Name
}

func (ons NavigationService) GetRoute(origin, destination DTO.Location, waypoints []DTO.Location) (NavigationResponse, error) {
	var navigationResponse NavigationResponse
	request := NewOSRMRouteRequest(
		origin,
		destination,
		waypoints,
	)
	url := generateRouteUrl(ons, request)
	responseData, err := Http.Get(url, ons.DefaultServiceClient)
	if err != nil {
		log.WithError(err).Error("error while http get request")
		return NavigationResponse{}, err
	}
	err = json.Unmarshal(responseData, &navigationResponse)
	if err != nil {
		log.WithError(err).Error("Error while unmarshalling the navigation response")
		return NavigationResponse{}, err
	}
	return navigationResponse, nil
}
