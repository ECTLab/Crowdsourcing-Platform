	package Route

	import (
		"navigation/pkg/Clients/Osrm"
		"navigation/pkg/Clients/Redis"
		"navigation/pkg/DTO"

		log "github.com/sirupsen/logrus"
	)

	func GetRoute(request *DTO.NavigationRequest) (Osrm.NavigationResponse, error) {
		osrmProvider := Osrm.NavigationOsrmService
		origin := request.Origin
		destination := request.Destination
		waypoints := request.MiddleDestinations

		response, err := osrmProvider.GetRoute(origin, destination, waypoints)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"response":       response,
				"origin":         origin,
				"destination":    destination,
				"waypoints":      waypoints,
			}).Error("error in calling route navigation", err)
			return response, err
		}
		addPoliceToResponse(&response, Redis.PoliceData)
		return response, nil
	}
