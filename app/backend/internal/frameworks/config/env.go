package config

type EnvironmentName = string

const (
	EnvironmentLocal      EnvironmentName = "local"
	EnvironmentTest       EnvironmentName = "test"
	EnvironmentStage      EnvironmentName = "stage"
	EnvironmentProduction EnvironmentName = "production"
)
