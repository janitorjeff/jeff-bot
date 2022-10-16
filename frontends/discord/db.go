package discord

import (
	"database/sql"
	"fmt"

	"git.slowtyper.com/slowtyper/janitorjeff/core"

	dg "github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func getScope(t int, id string, msg *dg.Message) (int64, error) {
	db := core.Globals.DB
	db.Lock.Lock()
	defer db.Lock.Unlock()

	switch t {
	case Default:
		return getScopeDefault(t, msg.ChannelID, msg.GuildID)
	case Guild, Channel, Thread:
		// TODO
		fallthrough
	case User:
		return getScopeUser(id)
	case Author:
		return getScopeUser(msg.Author.ID)
	default:
		return -1, fmt.Errorf("type '%d' not supported", t)
	}
}

func getScopeDefault(type_ int, channelID, guildID string) (int64, error) {
	// In some cases a guild does not exist, for example in a DM, thus we are
	// forced to use the channel scope. A guild id is also not included in the
	// message object returned after sending a message, only the channel id is.
	// So, it can be difficult to differentiate if a message comes from a DM
	// or from a returned message object, since in order to do so we rely on
	// checking if the guild id field is empty. A way to solve this is relying
	// on the fact that the returned message object only comes after a user
	// has executed a command in that scope. That means that *if* a guild
	// exists it will already have been added in the database, and so we use
	// that.
	//
	// This can break if for example the bot were to send a message in a scope
	// where no message has ever been sent which means that channel/guild ids
	// have not been recorded, since the message create hook ignores the bot's
	// messages. This means that with the current implementation if that
	// were to happen in a guild, then channel scoped would be returned instead
	// of the guild scope. This is not a problem in places where no guild
	// exists.

	switch type_ {
	case Default, Guild, Channel, Thread:
		break
	default:
		return -1, fmt.Errorf("type '%d' not supported", type_)
	}

	if type_ == Thread {
		return -1, fmt.Errorf("Thread scopes not supported yet")
	}

	// if scope exists return it instead of re-adding it
	channelScope, err := dbGetChannelScope(channelID)
	if err == nil {
		if type_ == Channel {
			return channelScope, nil
		}

		// find channel's guild scope
		// if channel exists then guild does also, even if it's the special
		// empty guild
		guildScope, err := dbGetGuildFromChannel(channelScope)
		// A guild does exist even if it's a DM, it's the empty string guild,
		// so this means if there's an error, it's a different kind.
		if err != nil {
			return -1, err
		}

		if type_ == Guild {
			return guildScope, nil
		}

		// In the schema we make it so that the empty guild is the first one
		// added, and has `scope = 1`. This is the only way I can come up with
		// that doesn't require reading the DB again to search for the guild's
		// id
		if guildScope == 1 {
			return channelScope, err
		}
		return guildScope, nil
	}

	db := core.Globals.DB

	tx, err := db.DB.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	// only create a new guildScope if it doesn't already exist
	guildScope, err := dbGetGuildScope(guildID)
	if err != nil {
		guildScope, err = dbAddGuildScope(tx, guildID)
		if err != nil {
			return -1, err
		}
	}

	channelScope, err = dbAddChannelScope(tx, channelID, guildScope)
	if err != nil {
		return -1, err
	}

	var scope int64

	switch type_ {
	case Guild:
		scope = guildScope
	case Channel:
		scope = channelScope
	default:
		// We are sure that no guild exists here, which is why the channel is
		// returned
		if guildID == "" {
			scope = channelScope
		} else {
			scope = guildScope
		}
	}

	return scope, tx.Commit()
}

func getScopeUser(id string) (int64, error) {
	scope, err := dbGetUserScope(id)
	if err == nil {
		return scope, nil
	}

	// TODO: check if id exists

	db := core.Globals.DB

	tx, err := db.DB.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	scope, err = dbAddUserScope(tx, id)
	if err != nil {
		return -1, err
	}

	return scope, tx.Commit()
}

func dbAddGuildScope(tx *sql.Tx, guildID string) (int64, error) {
	db := core.Globals.DB

	scope, err := db.ScopeAdd(tx)
	if err != nil {
		return -1, err
	}

	_, err = tx.Exec(`
		INSERT OR IGNORE INTO PlatformDiscordGuilds(id, guild)
		VALUES (?, ?)`, scope, guildID)

	if err != nil {
		return -1, err
	}

	return scope, nil
}

func dbAddChannelScope(tx *sql.Tx, channelID string, guildScope int64) (int64, error) {
	db := core.Globals.DB

	scope, err := db.ScopeAdd(tx)
	if err != nil {
		return -1, err
	}

	_, err = tx.Exec(`
		INSERT OR IGNORE INTO PlatformDiscordChannels(id, channel, guild)
		VALUES (?, ?, ?)`, scope, channelID, guildScope)

	if err != nil {
		return -1, err
	}

	return scope, nil
}

func dbGetGuildScope(guildID string) (int64, error) {
	db := core.Globals.DB

	row := db.DB.QueryRow(`
		SELECT id
		FROM PlatformDiscordGuilds
		WHERE guild = ?`, guildID)

	var id int64
	err := row.Scan(&id)
	return id, err
}

func dbGetChannelScope(channelID string) (int64, error) {
	db := core.Globals.DB

	row := db.DB.QueryRow(`
		SELECT id
		FROM PlatformDiscordChannels
		WHERE channel = ?`, channelID)

	var id int64
	err := row.Scan(&id)
	return id, err
}

func dbGetGuildFromChannel(channelScope int64) (int64, error) {
	db := core.Globals.DB

	row := db.DB.QueryRow(`
		SELECT guild
		FROM PlatformDiscordChannels
		WHERE id = ?`, channelScope)

	var guildScope int64
	err := row.Scan(&guildScope)
	return guildScope, err
}

func dbAddUserScope(tx *sql.Tx, userID string) (int64, error) {
	db := core.Globals.DB

	scope, err := db.ScopeAdd(tx)
	if err != nil {
		return -1, err
	}

	_, err = tx.Exec(`
		INSERT INTO PlatformDiscordUsers(id, user)
		VALUES (?, ?)`, scope, userID)

	log.Debug().
		Err(err).
		Int64("scope", scope).
		Str("user", userID).
		Msg("added user scope to db")

	return scope, err
}

func dbGetUserScope(userID string) (int64, error) {
	db := core.Globals.DB

	row := db.DB.QueryRow(`
		SELECT id
		FROM PlatformDiscordUsers
		WHERE user = ?`, userID)

	var id int64
	err := row.Scan(&id)

	log.Debug().
		Err(err).
		Int64("scope", id).
		Str("user", userID).
		Msg("got user scope from db")

	return id, err
}