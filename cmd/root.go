package cmd

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nint8835/discord-vc-stats/pkg/config"
)

var rootCmd = &cobra.Command{
	Use:   "discord-vc-stats",
	Short: "Discord bot that collects stats on voice channel activity",
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
