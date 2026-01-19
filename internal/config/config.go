package config

type Config struct {
	Service  *ServiceConfig  `json:"service"`
	Database *DatabaseConfig `json:"database"`
}
