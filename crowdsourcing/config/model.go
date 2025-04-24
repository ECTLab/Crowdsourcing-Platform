package config

type Config struct {
	LogLevel                       string                         `mapstructure:"log_level"`
	ServiceConfigs                 ServiceConfigs                 `mapstructure:"service"`
	AggregationConfigs             AggregationConfigs             `mapstructure:"aggregation"`
	RedisCrowdsourcing             RedisCrowdsourcing             `mapstructure:"redis_crowdsourcing"`
	RedisAnnotation                RedisAnnotation                `mapstructure:"redis_annotation"`
	OsrmTable                      OsrmTable                      `mapstructure:"osrm_table"`
	OsrmMatching                   OsrmMatching                   `mapstructure:"osrm_matching"`
	OnlineReportVisibilityOverHttp OnlineReportVisibilityOverHttp `mapstructure:"online_report_visibility_over_http"`
}

type AggregationConfigs struct {
	PrintPercentInterval     int     `mapstructure:"print_percent_interval"`
	ClusterDurationThreshold float64 `mapstructure:"cluster_duration_threshold"`
	ClusterLouvainResolution float64 `mapstructure:"cluster_louvain_resolution"`
	MaxOsrmTableLocations    int     `mapstructure:"max_osrm_table_locations"`
	OsmPbfUrl                string  `mapstructure:"osm_pbf_url"`
	OsmPbfVersionUrl         string  `mapstructure:"osm_pbf_version_url"`
	OsmUser                  string  `mapstructure:"osm_user"`
	OsmPass                  string  `mapstructure:"osm_pass"`
	OsmDirPath               string  `mapstructure:"osm_dir_path"`
	InsecureOSMDownload      bool    `mapstructure:"insecure_osm_download"`
}

type ServiceConfigs struct {
	ServiceProtoPort int32  `mapstructure:"service_proto_port"`
	ServiceProtoHost string `mapstructure:"service_proto_host"`
	ServiceName      string `mapstructure:"service_name"`
}


type RedisCrowdsourcing struct {
	Enabled                bool     `mapstructure:"enabled"`
	IsSentinel             bool     `mapstructure:"is_sentinel"`
	Addresses              []string `mapstructure:"addresses"`
	SentinelName           string   `mapstructure:"sentinel_name"`
	EventDefaultTtlMinutes int32    `mapstructure:"event_default_ttl_minutes"`
}

type RedisAnnotation struct {
	Enabled                bool     `mapstructure:"enabled"`
	IsSentinel             bool     `mapstructure:"is_sentinel"`
	Addresses              []string `mapstructure:"addresses"`
	SentinelName           string   `mapstructure:"sentinel_name"`
	EventDefaultTtlMinutes int32    `mapstructure:"event_default_ttl_minutes"`
}

type OsrmTable struct {
	Address   string `mapstructure:"address,omitempty"`
	TimeoutMS int32  `mapstructure:"timeout_ms,omitempty"`
}

type OsrmMatching struct {
	Address   string `mapstructure:"address,omitempty"`
	TimeoutMS int32  `mapstructure:"timeout_ms,omitempty"`
}

type OnlineReportVisibilityOverHttp struct {
	Enabled     bool   `mapstructure:"enabled"`
	AccessToken string `mapstructure:"access_token"`
}
