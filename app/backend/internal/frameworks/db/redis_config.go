package db

import "api/internal/frameworks/config"

type RedisMode = string

const (
	RedisModeSingle   RedisMode = "single"
	RedisModeSentinel RedisMode = "sentinel"
	RedisModeCluster  RedisMode = "cluster"
)

type RedisConfig struct {
	Active          bool      `mapstructure:"active"`
	Mode            RedisMode `mapstructure:"mode"`
	Username        string    `mapstructure:"username"`
	Password        string    `mapstructure:"password"`
	Endpoints       []string  `mapstructure:"endpoints"`
	DB              int       `mapstructure:"db"`
	MasterName      string    `mapstructure:"masterName"`
	MaxRetries      int       `mapstructure:"maxRetries"`
	MinRetryBackoff string    `mapstructure:"minRetryBackoff"`
	MaxRetryBackoff string    `mapstructure:"maxRetryBackoff"`
	PoolSize        int       `mapstructure:"poolSize"`
	MinIdleConns    int       `mapstructure:"minIdleConns"`
	ConnMaxLifetime string    `mapstructure:"connMaxLifetime"`
	PoolTimeout     string    `mapstructure:"poolTimeout"`
	ConnMaxIdleTime string    `mapstructure:"connMaxIdleTime"`
}

func NewRedisConfig(configHelper config.ConfigHelper) (*RedisConfig, error) {
	var cfg RedisConfig
	if err := configHelper.UnmarshalKey("components.db.redis", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
