package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nint8835/discord-vc-stats/pkg/bot"
	"github.com/nint8835/discord-vc-stats/pkg/metrics"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the bot",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		botInst, err := bot.New()
		if err != nil {
			log.Fatal().Err(err).Msg("Error creating bot")
		}

		go func() {
			err = botInst.Start()
			if err != nil {
				log.Fatal().Err(err).Msg("Error starting bot")
			}
		}()

		err = metrics.ServeMetrics()
		if err != nil {
			log.Fatal().Err(err).Msg("Error serving metrics")
		}

		botInst.Stop()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
