package config

import "go.uber.org/zap"

const (
	DefaultBuildVersion string = "development"
	DefaultAppName      string = "app"
)

/*
Application creation when creating a package body, you can use ldflags to change it.
Please refer to https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
*/
var (
	//app name
	appName string = DefaultAppName
	//environment name of app
	envName EnvironmentName = EnvironmentLocal
	//build version
	buildVersion string = DefaultBuildVersion
	//build time
	buildTime string
	//git commit hash
	buildCommit string
	//git branch
	buildBranch string
)

// Get the program name
func GetAppName() string {
	return appName
}

// Get the application environment
func GetEnvName() EnvironmentName {
	return envName
}

// Get the application creation version number
func GetBuildVersion() string {
	return buildVersion
}

// Get the application creation time
func GetBuildTime() string {
	return buildTime
}

// Get the commit hash of the application git
func GetBuildCommit() string {
	return buildCommit
}

// Get the application git branch
func GetBuildBranch() string {
	return buildBranch
}

func IsDebugging() bool {
	return buildVersion == DefaultBuildVersion
}

func PrintAppInfo(logger *zap.Logger) {

	logger.Info("App Information",
		zap.String("app-name", appName),
		zap.String("env-name", envName),
		zap.String("build-version", buildVersion),
		zap.String("build-time", buildTime),
		zap.String("build-commit", buildCommit),
		zap.String("build-branch", buildBranch))
}
