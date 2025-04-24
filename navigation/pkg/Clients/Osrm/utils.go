package Osrm

import (
	"fmt"
	"navigation/config"
	"navigation/tools/util"
	"strings"
)



func generateRouteUrl(ons NavigationService, request *RouteRequest) string {
	baseURL := fmt.Sprintf("%s/%s", ons.Host, config.GetServiceConfig().Osrm.Route.Url)
	originStr := util.ConvertLocationToString(request.Origin)
	destinationStr := util.ConvertLocationToString(request.Destination)
	waypointsStr := ""
	coords := fmt.Sprintf("%s;%s", originStr, destinationStr)

	if len(request.Waypoints) > 0 {
		waypoints := util.ConvertLocationsToString(request.Waypoints)
		waypointsStr = strings.Join(waypoints, ";")
		coords = fmt.Sprintf("%s;%s;%s", originStr, waypointsStr, destinationStr)
	}

	urlParams := []string{
		"overview=full",
		"steps=true",
		"annotations=true",
		fmt.Sprintf("alternatives=%v", request.Alternative),
	}

	url := fmt.Sprintf("%s/%s?%s", baseURL, coords, strings.Join(urlParams, "&"))

	return url
}
