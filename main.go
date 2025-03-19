package main

import "C"
import (
	"blog/configs"
	"blog/internal/platform/pkg/lang"
	"blog/internal/ui/cli"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

var mainCommand = &cobra.Command{}

func main() {
	mainCommand.AddCommand(cli.Run(config))
	mainCommand.AddCommand(cli.Seeder(config))
	err := mainCommand.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
