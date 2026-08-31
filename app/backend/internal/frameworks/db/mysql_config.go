package db

import (
	"api/internal/frameworks/config"
)

type MySQLConfig struct {
	Active          bool   `mapstructure:"active"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	Charset         string `mapstructure:"charset"`
	ParseTime       bool   `mapstructure:"parseTime"`
	Loc             string `mapstructure:"loc"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`
	ConnMaxLifetime string `mapstructure:"connMaxLifetime"`
	ConnMaxIdleTime string `mapstructure:"connMaxIdleTime"`
	Timeout         string `mapstructure:"timeout"`
	ReadTimeout     string `mapstructure:"readTimeout"`
	WriteTimeout    string `mapstructure:"writeTimeout"`
	Debug           bool   `mapstructure:"debug"`
}

func NewMysqlConfig(configHelper config.ConfigHelper) (*MySQLConfig, error) {
	var cfg MySQLConfig
	if err := configHelper.UnmarshalKey("components.db.mysql", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
