package config

import (
	"os"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

const (
	DefaultConfigDir  string = "."
	defaultConfigName string = "config"
)

var (
	CommonViperDecoder viper.DecoderConfigOption = viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			mapstructure.DecodeHookFunc(ParseEnvViperDecodeHook),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		))
)

// Parse environment variables ${env_xxx} or ${env_xxx:default value}
func ParseEnvViperDecodeHook(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() == reflect.String {
		stringData := data.(string)

		if strings.HasPrefix(stringData, "${") && strings.HasSuffix(stringData, "}") {
			i := strings.Index(stringData, ":")
			prefix := strings.Index(stringData, "${") + 2
			suffix := strings.Index(stringData, "}")
			def := ""
			envKey := ""

			if i > -1 {
				def = stringData[i+1 : suffix]
				envKey = stringData[prefix:i]
			} else {
				envKey = stringData[prefix:suffix]
			}

			envVarValue := os.Getenv(envKey)
			if envVarValue == "" {
				envVarValue = def
			}

			return envVarValue, nil
		}
	}
	return data, nil
}

func NewConfigNames(major string, minors ...string) []string {
	results := make([]string, len(minors)+1)
	results[0] = major
	for i, minor := range minors {
		results[i+1] = major + "." + minor
	}
	return results
}

func Load(dir string, names ...string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(dir)
	v.SetConfigType("yaml")
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	if len(names) == 0 {
		names = append(names, defaultConfigName)
	}

	for i, n := range names {

		v.SetConfigName(n)
		// v.SetConfigFile("./configs/config.local.yaml")

		if i == 0 {
			if err := v.ReadInConfig(); err != nil {
				return nil, err
			}
		} else {
			v.MergeInConfig()
		}
	}

	return v, nil
}
