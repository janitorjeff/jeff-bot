package discord

import (
	"fmt"
	"sync"

	"git.slowtyper.com/slowtyper/janitorjeff/core"

	dg "github.com/bwmarrin/discordgo"
)

var replies = replyCache{}

type replyCache struct {
	lock    sync.RWMutex
	replies map[string]string
}

func (r *replyCache) Set(key, value string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.replies == nil {
		r.replies = map[string]string{}
	}
	r.replies[key] = value
}

func (r *replyCache) Get(key string) (string, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	value, ok := r.replies[key]
	return value, ok
}

type DiscordMessageEdit struct {
	Session *dg.Session
	Message *dg.MessageUpdate
}

func (d *DiscordMessageEdit) Admin() bool {
	return isAdmin(d.Message.Author.ID)
}

func (d *DiscordMessageEdit) Parse() (*core.Message, error) {
	msg := parse(d.Message.Message)
	msg.Client = d
	return msg, nil
}

func (d *DiscordMessageEdit) Scope(type_ int) (int64, error) {
	return getScope(type_, d.Message.ChannelID, d.Message.GuildID, d.Message.Author.ID)
}

func (d *DiscordMessageEdit) Write(msg interface{}, usrErr error) (*core.Message, error) {
	switch t := msg.(type) {
	case string:
		return sendText(d.Session, msg.(string), d.Message.ChannelID)

	case *dg.MessageEmbed:
		embed := msg.(*dg.MessageEmbed)
		id, ok := replies.Get(d.Message.ID)
		if !ok {
			return sendEmbed(d.Session, d.Message.Message, embed, usrErr)
		}
		return editEmbed(d.Session, embed, usrErr, id, d.Message.ChannelID)

	default:
		return nil, fmt.Errorf("Can't send discord message of type %v", t)
	}

}

func messageEdit(s *dg.Session, m *dg.MessageUpdate) {
	if m == nil {
		fmt.Println("MESSAGE IS NIL FOR SOME REASON")
		fmt.Println("MESSAGE IS NIL FOR SOME REASON")
		fmt.Println("MESSAGE IS NIL FOR SOME REASON")
		fmt.Println("MESSAGE IS NIL FOR SOME REASON")
		fmt.Println("MESSAGE IS NIL FOR SOME REASON")
		fmt.Println("MESSAGE IS NIL FOR SOME REASON")
		fmt.Println("MESSAGE IS NIL FOR SOME REASON")
	}

	if m.Author.Bot {
		return
	}

	if len(m.Content) == 0 {
		return
	}

	d := &DiscordMessageEdit{s, m}
	msg, err := d.Parse()
	if err != nil {
		return
	}
	msg.Run()
}