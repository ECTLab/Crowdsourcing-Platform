package Route

import (
	"math/rand/v2"
	"navigation/pkg/Clients/Osrm"
	"navigation/pkg/Clients/Redis"
	"navigation/pkg/DTO"
)


func addPoliceToResponse(response *Osrm.NavigationResponse, policeData map[Redis.PoliceRedisKey]Redis.OnlineReportRedisSchema) {
	for i, route := range response.Routes {
		for j, leg := range route.Legs {
			polices := getPolices(leg.Annotation.Nodes, policeData)
			for i := range polices {
				if polices[i].Confidence < rand.Float32() {
					polices[i] = DTO.Police{Exists: false}
				}
			}
			response.Routes[i].Legs[j].Annotation.Police = polices
		}
	}
}


func getPolices(nodes []int64, policeData map[Redis.PoliceRedisKey]Redis.OnlineReportRedisSchema) []DTO.Police {
	polices := []DTO.Police{
		{Exists: false}, // client can not show police on first segment
	}
	for n := 1; n < len(nodes)-1; n++ {
		u := nodes[n]
		v := nodes[n+1]

		routeIsFromVToU := false

		if v > u {
			t := v
			v = u
			u = t
			routeIsFromVToU = true
		}

		police, exists := policeData[Redis.PoliceRedisKey{U: u, V: v}]

		if exists {
			var offset float64
			if routeIsFromVToU {
				offset = police.OffsetRest
			} else {
				offset = police.Offset
			}

			polices = append(polices, DTO.Police{
				Id:           police.Uid,
				Exists:       true,
				Offset:       offset,
				Confirmation: police.Confirmation,
				Confidence:   police.Confidence,
			})
		} else {
			polices = append(polices, DTO.Police{Exists: false})
		}
	}
	return polices
}


