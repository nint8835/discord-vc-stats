package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"

	"github.com/nint8835/discord-vc-stats/pkg/config"
	"github.com/nint8835/discord-vc-stats/pkg/metrics"
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

func handleVoiceStateUpdate(session *discordgo.Session, update *discordgo.VoiceStateUpdate) {
	user, err := session.User(update.UserID)
	if err != nil {
		log.Error().Err(err).Str("user_id", update.UserID).Msg("Error getting member")
		return
	}

	// User was not previously in a voice channel & joined a voice channel
	if update.ChannelID != "" && (update.BeforeUpdate == nil || update.BeforeUpdate.ChannelID == "") {
		channel, err := session.Channel(update.ChannelID)
		if err != nil {
			log.Error().Err(err).Str("channel_id", update.ChannelID).Msg("Error getting channel")
			return
		}
		metrics.MembersInVoice.WithLabelValues(channel.Name, user.Username).Set(1)
	}
	// User moved from one voice channel to another
	if update.ChannelID != "" && (update.BeforeUpdate != nil && update.BeforeUpdate.ChannelID != "") && update.ChannelID != update.BeforeUpdate.ChannelID {
		newChannel, err := session.Channel(update.ChannelID)
		if err != nil {
			log.Error().Err(err).Str("channel_id", update.ChannelID).Msg("Error getting channel")
			return
		}
		oldChannel, err := session.Channel(update.BeforeUpdate.ChannelID)
		if err != nil {
			log.Error().Err(err).Str("channel_id", update.BeforeUpdate.ChannelID).Msg("Error getting channel")
			return
		}

		metrics.MembersInVoice.WithLabelValues(oldChannel.Name, user.Username).Set(0)
		metrics.MembersInVoice.WithLabelValues(newChannel.Name, user.Username).Set(1)
	}
	// User left voice
	if update.ChannelID == "" && (update.BeforeUpdate != nil && update.BeforeUpdate.ChannelID != "") {
		channel, err := session.Channel(update.BeforeUpdate.ChannelID)
		if err != nil {
			log.Error().Err(err).Str("channel_id", update.BeforeUpdate.ChannelID).Msg("Error getting channel")
			return
		}
		metrics.MembersInVoice.WithLabelValues(channel.Name, user.Username).Set(0)
	}
}

func New() (*Bot, error) {
	bot := &Bot{
		quitChan: make(chan struct{}),
	}

	session, err := discordgo.New("Bot " + config.Instance.DiscordToken)
	if err != nil {
		return nil, fmt.Errorf("error creating discord session: %w", err)
	}
	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsMessageContent
	bot.Session = session

	session.AddHandler(handleVoiceStateUpdate)

	return bot, nil
}
