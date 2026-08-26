package db

import "api/internal/frameworks/config"

type MongoConfig struct {
	Active                 bool     `mapstructure:"active"`
	Hosts                  []string `mapstructure:"hosts"`
	ReplicaSet             string   `mapstructure:"replicaSet"`
	MaxPoolSize            uint64   `mapstructure:"maxPoolSize"`
	MinPoolSize            uint64   `mapstructure:"minPoolSize"`
	ConnectTimeoutMS       int      `mapstructure:"connectTimeoutMS"`
	ServerSelectionTimeout int      `mapstructure:"serverSelectionTimeoutMS"`
	HeartbeatIntervalMS    int      `mapstructure:"heartbeatIntervalMS"`
	//option: primary, primaryPreferred, secondary, secondaryPreferred, nearest
	ReadPreference string `mapstructure:"readPreference"`
	RetryWrites    bool   `mapstructure:"retryWrites"`
	RetryReads     bool   `mapstructure:"retryReads"`
	// e.g. majority, w:1, w:2, w:3, w: "majority", w: "tagSet", wtimeout: 1000
	WriteConcern string         `mapstructure:"writeConcern"`
	Username     string         `mapstructure:"username"`
	Password     string         `mapstructure:"password"`
	DbName       string         `mapstructure:"dbname"`
	Log          MongoLogConfig `mapstructure:"log"`
}

type MongoLogConfig struct {
	Enabled bool `mapstructure:"enabled"`
	Debug   bool `mapstructure:"debug"`
}

func NewMongoConfig(configHelper config.ConfigHelper) (*MongoConfig, error) {
	var cfg MongoConfig
	if err := configHelper.UnmarshalKey("components.db.mongodb", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
