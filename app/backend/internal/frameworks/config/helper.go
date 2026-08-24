package config

import (
	"github.com/cockroachdb/errors"
	"github.com/spf13/viper"
)

type ConfigHelper interface {
	Unmarshal(val any) error
	UnmarshalKey(key string, val any) error
}
type configHelper struct {
	viper *viper.Viper
}

func (helper *configHelper) Unmarshal(val any) error {
	return errors.WithStack(helper.viper.Unmarshal(val, CommonViperDecoder))
}
func (helper *configHelper) UnmarshalKey(key string, val any) error {
	return errors.WithStack(helper.viper.UnmarshalKey(key, val, CommonViperDecoder))
}

func NewConfigHelper(configDir string, configNames ...string) (ConfigHelper, error) {
	if len(configNames) == 0 {
		return nil, errors.New("configNames empty")
	}
	viper, err := Load(
		configDir,
		NewConfigNames(configNames[0], configNames[1:]...)...)
	if err != nil {
		return nil, err
	}
	return &configHelper{viper: viper}, nil
}
