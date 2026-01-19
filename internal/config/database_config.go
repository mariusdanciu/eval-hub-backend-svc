package config

type DatabaseConfig struct {
	URL      string `mapstructure:"url,omitempty"`
	User     string `mapstructure:"user,omitempty"`
	Password string `mapstructure:"password,omitempty"`
	Name     string `mapstructure:"name,omitempty"`
	Port     string `mapstructure:"port,omitempty"`
	Host     string `mapstructure:"host,omitempty"`
	Schema   string `mapstructure:"schema,omitempty"`
	SSLMode  string `mapstructure:"ssl_mode,omitempty"`
}
