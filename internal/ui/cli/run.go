package cli

import (
	"blog/configs"
	"blog/internal/platform/application"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Run(config *configs.Config) *cobra.Command {
	return &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := application.NewApp(config)
			if err != nil {
				log.WithError(err).Fatal("failed to create instance from application")
				return err
			}
			app.RunRouter()
			return app.Close()
		},
	}
}
