package discord

import (
	"fmt"

	"git.slowtyper.com/slowtyper/janitorjeff/core"

	dg "github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type DiscordMessageCreate struct {
	Session *dg.Session
	Message *dg.MessageCreate
}

func (d *DiscordMessageCreate) Parse() (*core.Message, error) {
	msg := parse(d.Message.Message)
	msg.Client = d
	return msg, nil
}

func (d *DiscordMessageCreate) Scope(type_ int) (int64, error) {
	return getScope(type_, d.Message.ChannelID, d.Message.GuildID, d.Message.Author.ID)
}

func (d *DiscordMessageCreate) Write(msg interface{}, usrErr error) (*core.Message, error) {
	switch t := msg.(type) {
	case string:
		text := msg.(string)
		lenLim := 2000
		// TODO: grapheme clusters instead of plain len?
		lenCnt := func(s string) int { return len(s) }
		return messagesTextSend(d.Session, text, d.Message.ChannelID, lenLim, lenCnt)

	case *dg.MessageEmbed:
		// TODO: implement message scrolling
		embed := msg.(*dg.MessageEmbed)
		if embed.Color == 0 {
			// default value of EmbedColor is 0 so even if it's not been set
			// then everything should be ok
			if usrErr == nil {
				embed.Color = core.Globals.Discord.EmbedColor
			} else {
				embed.Color = core.Globals.Discord.EmbedErrColor
			}
		}

		// TODO: Consider adding an option which allows one of these 3 values
		// - no reply + no ping, just an embed
		// - reply + no ping (default)
		// - reply + ping
		// Maybe even no embed and just plain text?
		m := &dg.MessageSend{
			Embeds: []*dg.MessageEmbed{
				embed,
			},
			AllowedMentions: &dg.MessageAllowedMentions{
				Parse: []dg.AllowedMentionType{}, // don't ping user
			},
			Reference: d.Message.Reference(),
		}

		// TODO: return message object
		_, err := d.Session.ChannelMessageSendComplex(d.Message.ChannelID, m)
		return nil, err

	default:
		return nil, fmt.Errorf("Can't send discord message of type %v", t)
	}

}

func messageCreate(s *dg.Session, m *dg.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Author.Bot {
		return
	}

	if len(m.Content) == 0 {
		return
	}

	d := &DiscordMessageCreate{s, m}
	msg, err := d.Parse()
	if err != nil {
		log.Debug().Err(err).Send()
		return
	}

	msg.Run()
}
