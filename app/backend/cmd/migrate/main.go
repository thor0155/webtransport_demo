package main

import (
	"api/internal/frameworks/config"

	"go.uber.org/zap"
)

func main() {

	configHelper, err := config.NewConfigHelper("./configs", "migrate", config.GetEnvName())
	if err != nil {
		panic(err)
	}

	log, err := GenLogger(configHelper)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	config.PrintAppInfo(log)

	lifecycle, err := GenLifecycle(log, configHelper)
	if err != nil {
		log.Error("gen lifecycle", zap.Error(err))
		return
	}

	if err := lifecycle.Run(); err != nil {
		log.Fatal("lifecycle stopped with error", zap.Error(err))
	}
}
