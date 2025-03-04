package main

import (
	"blog/configs"
	"blog/internal/platform/application"
	log "github.com/sirupsen/logrus"
)

var (
	config *configs.Config
)

func init() {
	cfg, err := configs.Load("./configs/config.yaml")
	if err != nil {
		log.WithError(err).Fatal("failed to load configurations")
	}
	config = cfg
	logLevel, err := log.ParseLevel(cfg.App.LogLevel)
	if err != nil {
		log.WithError(err).Fatal("failed to parse log level")
	}
	log.SetLevel(logLevel)
}
func main() {
	app := application.NewApp(config)
	app.RunRouter()
}
