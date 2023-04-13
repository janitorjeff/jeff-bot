package discord

import (
	"fmt"
	"io"
	"strings"

	"github.com/janitorjeff/jeff-bot/core"

	dg "github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type MessageCreate struct {
	Session *dg.Session
	Message *dg.MessageCreate
	VC      *dg.VoiceConnection
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

	// TODO: remove this when each server can configure which commands will be
	// active
	if m.GuildID == "348368013382254602" && strings.HasPrefix(m.Content, "!") {
		if !strings.HasPrefix(m.Content, "!pb") {
			return
		}
	}

	d := &MessageCreate{
		Session: s,
		Message: m,
	}
	msg, err := d.Parse()
	if err != nil {
		log.Debug().Err(err).Send()
		return
	}

	msg.Run()
}

///////////////
//           //
// Messenger //
//           //
///////////////

func (d *MessageCreate) Parse() (*core.Message, error) {
	msg := parse(d.Message.Message)
	msg.Client = d
	msg.Speaker = d
	return msg, nil
}

func (d *MessageCreate) PersonID(s, placeID string) (string, error) {
	return getPersonID(s, placeID, d.Message.Author.ID, d.Session)
}

func (d *MessageCreate) PlaceID(s string) (string, error) {
	return getPlaceID(s, d.Session)
}

func (d *MessageCreate) Person(id string) (int64, error) {
	return getPersonScope(id)
}

func (d *MessageCreate) PlaceExact(id string) (int64, error) {
	return getPlaceExactScope(id, d.Message.ChannelID, d.Message.GuildID, d.Session)
}

func (d *MessageCreate) PlaceLogical(id string) (int64, error) {
	return getPlaceLogicalScope(id, d.Message.ChannelID, d.Message.GuildID, d.Session)
}

func (d *MessageCreate) Usage(usage string) any {
	return getUsage(usage)
}

func (d *MessageCreate) send(msg any, usrErr error, ping bool) (*core.Message, error) {
	switch t := msg.(type) {
	case string:
		return sendText(d.Session, d.Message.Message, msg.(string), ping)
	case *dg.MessageEmbed:
		embed := msg.(*dg.MessageEmbed)
		return sendEmbed(d.Session, d.Message.Message, embed, usrErr, ping)
	default:
		return nil, fmt.Errorf("Can't send discord message of type %v", t)
	}
}

func (d *MessageCreate) Send(msg any, usrErr error) (*core.Message, error) {
	return d.send(msg, usrErr, false)
}

func (d *MessageCreate) Ping(msg any, usrErr error) (*core.Message, error) {
	return d.send(msg, usrErr, true)
}

func (d *MessageCreate) Write(msg any, usrErr error) (*core.Message, error) {
	return d.Send(msg, usrErr)
}

/////////////
//         //
// Speaker //
//         //
/////////////

func (d *MessageCreate) Voice() bool {
	return true
}

func (d *MessageCreate) FrameRate() int {
	return frameRate
}

func (d *MessageCreate) Channels() int {
	return channels
}

func (d *MessageCreate) Join() error {
	v, err := joinUserVoiceChannel(d.Session, d.Message.GuildID, d.Message.Author.ID)
	if err != nil {
		return err
	}
	d.VC = v
	return nil
}

func (d *MessageCreate) Say(buf io.Reader, s *core.State) error {
	return voicePlay(d.VC, buf, s)
}
