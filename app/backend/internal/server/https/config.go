package https

import "api/internal/frameworks/config"

type HttpServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Cors CORS   `mapstructure:"cors"`
}

type CORS struct {
	Enabled          bool
	AllowCredentials bool
	AllowOrigins     []string
	ExposeHeaders    []string
}

func ProvideHttpServerConfig(helper config.ConfigHelper) (*HttpServerConfig, error) {
	var r HttpServerConfig
	err := helper.UnmarshalKey("servers.http-server", &r)
	return &r, err
}
