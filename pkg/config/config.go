package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	LogLevel string `split_words:"true" default:"info"`

	BindAddr string `default:":12500" split_words:"true"`

	DiscordToken string `split_words:"true" required:"true"`
}

var Instance *Config

func Load() error {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Warn().Err(err).Msg("error loading .env file")
	}

	var newConfig Config
	err = envconfig.Process("vc_stats", &newConfig)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	logLevel, err := zerolog.ParseLevel(newConfig.LogLevel)
	if err != nil {
		return fmt.Errorf("error parsing provided log level: %w", err)
	}
	zerolog.SetGlobalLevel(logLevel)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	Instance = &newConfig

	return nil
}
