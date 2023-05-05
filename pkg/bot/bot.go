package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"

	"github.com/nint8835/discord-vc-stats/pkg/config"
)

type Bot struct {
	Session  *discordgo.Session
	quitChan chan struct{}
}

func (b *Bot) Start() error {
	log.Info().Msg("Starting bot...")

	err := b.Session.Open()
	if err != nil {
		return fmt.Errorf("error opening discord session: %w", err)
	}

	<-b.quitChan

	err = b.Session.Close()
	if err != nil {
		return fmt.Errorf("error closing discord session: %w", err)
	}

	return nil
}

func (b *Bot) Stop() {
	b.quitChan <- struct{}{}
}

func New() (*Bot, error) {
	bot := &Bot{
		quitChan: make(chan struct{}),
	}

	session, err := discordgo.New("Bot " + config.Instance.DiscordToken)
	if err != nil {
		return nil, fmt.Errorf("error creating discord session: %w", err)
	}
	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	bot.Session = session

	return bot, nil
}
