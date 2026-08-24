package wts

import (
	"api/internal/frameworks/config"
)

type WebTransportServerConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	CertFile    string `mapstructure:"cert-file"`
	KeyFile     string `mapstructure:"key-file"`
	Path        string `mapstructure:"path"`
	UseSelfCert bool   `mapstructure:"use-self-cert"`
}

func ProvideWebtransportServerConfig(helper config.ConfigHelper) (*WebTransportServerConfig, error) {
	var r WebTransportServerConfig
	err := helper.UnmarshalKey("servers.webtransport-server", &r)
	return &r, err
}
