package DTO

type TTSReply struct {
	RequestId string `json:"request_id"`
	Tid       string `json:"tid"`
	URL       string `json:"url"`
}

type NavigationResponse struct {
	RequestId string  `json:"request_id"`
	Routes    []Route `json:"routes"`
	TTL       int64   `json:"ttl"`
}

type GetVoiceLinkResponse struct {
	Address             string
	AddressAnnouncement string
	AnnouncementType    AnnouncementTypeEnum
}

type GetNearestStreetLocationResponse struct {
	Locations     []Location `json:"locations"`
}