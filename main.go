package main

import (
	"blog/configs"
	"blog/internal/platform/application"
	"blog/internal/platform/pkg/lang"
	"fmt"
	log "github.com/sirupsen/logrus"
	"runtime"
	"time"
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
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		PrettyPrint:     false,
		DataKey:         "data",
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", f.File, f.Line)
		},
	})
	lang.SetBundle()
	lang.SetDefaultLocale("en")
	lang.SetTranslationFilesPath("./assets/translations")
}
func main() {
	app, err := application.NewApp(config)
	if err != nil {
		log.WithError(err).Fatal("failed to create instance from application")
	}
	app.RunRouter()
}
