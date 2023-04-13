package discord

import (
	"fmt"
	"strings"

	"git.slowtyper.com/slowtyper/janitorjeff/core"
	"git.slowtyper.com/slowtyper/janitorjeff/utils"

	dg "github.com/bwmarrin/discordgo"
)

func getDisplayName(member *dg.Member, author *dg.User) string {
	var displayName string

	if member == nil || member.Nick == "" {
		displayName = author.Username
	} else {
		displayName = member.Nick
	}

	return displayName
}

func messagesTextSend(d *dg.Session, text, channel string, lenLim int, lenCnt func(string) int) (*core.Message, error) {
	var msg *dg.Message
	var err error

	if lenLim > lenCnt(text) {
		msg, err = d.ChannelMessageSend(channel, text)
	} else {
		parts := utils.Split(text, lenCnt, lenLim)
		for _, p := range parts {
			msg, err = d.ChannelMessageSend(channel, p)
		}
	}

	if err != nil {
		return nil, err
	}
	m := &DiscordMessage{d, msg}
	return m.Parse()
}

func PlaceInBackticks(s string) string {
	if !strings.Contains(s, "`") {
		return fmt.Sprintf("`%s`", s)
	}

	// Only way I could find to display backticks correctly. Very hacky.
	// Works for an arbitrary number of backticks. Works everywhere except
	// on android.

	const zeroWidthSpace = "\u200b"
	// zeroWidthSpace := "\u3164"

	s = strings.ReplaceAll(s, "`", zeroWidthSpace+"`"+zeroWidthSpace)
	return fmt.Sprintf("``%s``", s)
}