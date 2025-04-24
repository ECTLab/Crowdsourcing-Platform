package config

type Config struct {
	LogLevel          string            `mapstructure:"log_level"`
	ServiceConfigs    ServiceConfigs    `mapstructure:"service"`
	PrometheusConfigs PrometheusConfigs `mapstructure:"prometheus_configs"`
	KafkaConfigs      KafkaConfigs      `mapstructure:"kafka"`
	RedisAnnotation   Redis             `mapstructure:"redis_annotation"`
	RedisTTS          Redis             `mapstructure:"redis_tts"`
	Osrm              OsrmConfigs       `mapstructure:"osrm_configs"`
	Mongo             MongoConfigs      `mapstructure:"mongo"`
	ABTest            ABTestConfigs     `mapstructure:"ab_test"`
	Joob              Joob              `mapstructure:"joob"`
	NavigationOptions NavigationOptions `mapstructure:"navigation_options"`
	TTS               TTS               `mapstructure:"tts"`
	Annotation        Annotation        `mapstructure:"annotation"`
	Health            Health            `mapstructure:"health"`
	Product           Product           `mapstructure:"product"`
	DornaCsat         DornaCsatConfigs  `mapstructure:"dorna_csat"`
	Vroom			  Vroom 			`mapstructure:"vroom"`
}

type ServiceConfigs struct {
	ServiceProtoPort int32  `mapstructure:"service_proto_port"`
	ServiceProtoHost string `mapstructure:"service_proto_host"`
	Name             string `mapstructure:"name"`
}

type PrometheusConfigs struct {
	Host              string `mapstructure:"host"`
	Port              int32  `mapstructure:"port"`
	PushServerAddress string `mapstructure:"push_server_address"`
	MetricsPrefix     string `mapstructure:"metrics_prefix"`
	MetricsSuffix     string `mapstructure:"metrics_suffix"`
}

type KafkaConfigs struct {
	Enabled        bool     `mapstructure:"enabled"`
	Addresses      []string `mapstructure:"addresses"`
	DefaultVersion string   `mapstructure:"default_version"`
	Timeout        int      `mapstructure:"timeout"`
}

type Redis struct {
	Enabled                bool     `mapstructure:"enabled"`
	IsSentinel             bool     `mapstructure:"is_sentinel"`
	Addresses              []string `mapstructure:"addresses"`
	SentinelName           string   `mapstructure:"sentinel_name"`
	EventDefaultTtlMinutes int32    `mapstructure:"event_default_ttl_minutes"`
	ConnPoolMaxIdleTime    int      `mapstructure:"conn_pool_max_idle_time"`
}

type Osrm struct {
	Address   string `mapstructure:"address"`
	TimeoutMS int32  `mapstructure:"timeout_ms"`
}

type Vroom struct {
	Enabled		bool		`mapstructre:"enabled"`
	Address		string	    `mapstructure:"address"`
	Timeout 	int64		`mapstructure:"time_out"`
}

type MongoConfigs struct {
	Enabled                 bool              `mapstructure:"enabled"`
	ModelType               string            `mapstructure:"model_type"`
	Host                    string            `mapstructure:"host"`
	Port                    string            `mapstructure:"port"`
	Timeout                 string            `mapstructure:"timeout"`
	DbName                  string            `mapstructure:"db_name"`
	ModelVersionsCollection string            `mapstructure:"model_versions_collection"`
	User                    string            `mapstructure:"user"`
	Password                string            `mapstructure:"password"`
	Options                 map[string]string `mapstructure:"options"`
}

type OsrmConfigs struct {
	OsrmControl   Osrm            `mapstructure:"osrm_control"`
	OsrmTreatment Osrm            `mapstructure:"osrm_treatment"`
	Route         RouteConfigs    `mapstructure:"route"`
	Nearest       NearestConfigs  `mapstructure:"nearest"`
	Matching      MatchingConfigs `mapstructure:"matching"`
	Tsp 		  TspConfigs	  `mapstructure:"tsp"`
}

type TspConfigs struct {
	Url 			string 	`mapstructure:"url"`
	Address 		string 	`mapstructure:"address"`
}

type RouteConfigs struct {
	Url string `mapstructure:"url"`
}

type MatchingConfigs struct {
	Url       string `mapstructure:"url"`
	Address   string `mapstructure:"address"`
	MaxRadius string `mapstructure:"max_radius"`
}

type NearestConfigs struct {
	MaxValidZoomLevel float32  `mapstructure:"max_valid_zoom_level"`
	MinValidZoomLevel float32  `mapstructure:"min_valid_zoom_level"`
	SelectedCities    []string `mapstructure:"selected_cities"`
	Url               string   `mapstructure:"url"`
	Address           string   `mapstructure:"address"`
}

type ABTestConfigs struct {
	ABTestIntervalDurationMinutes int                       `mapstructure:"interval_duration_minutes"`
	DriverBased                   DriverBasedABTest         `mapstructure:"driver_based"`
	CityBased                     CityBasedABTest           `mapstructure:"city_based"`
}

type DornaCsatConfigs struct {
	TripPercentage uint32 `mapstructure:"trip_percentage"`
}


type CityBasedABTest struct {
	Enabled bool     `mapstructure:"enabled"`
	Cities  []string `mapstructure:"cities"`
}
type DriverBasedABTest struct {
	Enabled    		bool 	`mapstructure:"enabled"`
	ModNumber	 	int  	`mapstructure:"mod_number"`
	TreatmentMods	[]int	`mapstructure:"treatment_mods"`
}

type Joob struct {
	Enabled bool     `mapstructure:"enabled"`
	Address []string `mapstructure:"addresses"`
}

type NavigationOptions struct {
	NavigationTtlSeconds int64   `mapstructure:"ttl_seconds"`
	PickupRouteEnabled   bool    `mapstructure:"produce_navigation_response_on_joob_enabled"`
	TimeOffset           float64 `mapstructure:"time_offset"`
}

type TTS struct {
	Host                              string `mapstructure:"service_host"`
	Port                              int    `mapstructure:"service_port_grpc_proto"`
	Timeout                           int    `mapstructure:"timeout"`
	ProvideOnlyAltAnnouncementAddress bool   `mapstructure:"provide_only_alt_announcement_address"`
	ProduceNewSentenceOnJoobEnabled   bool   `mapstructure:"produce_new_sentence_on_joob_enabled"`
	AltAnnouncementTTL                int    `mapstructure:"alt_announcement_ttl"`
}

type Annotation struct {
	TrafficEnabled       		bool 						`mapstructure:"traffic_enabled"`
	EventsEnabled        		bool 						`mapstructure:"events_enabled"`
	OnlineReportsEnabled 		bool 						`mapstructure:"online_reports_enabled"`
	FeasibilityConfirmation		FeasibilityConfirmation		`mapstructure:"feasibility_confirmation"`
}

type FeasibilityConfirmation struct {
	Enabled					bool	`mapstructure:"enabled"`
	MaximumCasesInRide		int		`mapstructure:"maximum_cases_in_ride"`
	MaximumContributions	int		`mapstructure:"maximum_contributions"`
	NewDataFetchInterval    int     `mapstructure:"new_data_fetch_interval_min"`
}

type Health struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type Product struct {
	Name string `mapstructure:"name"`
}
