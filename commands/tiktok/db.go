package tiktok

import (
	"github.com/janitorjeff/jeff-bot/core"

	"github.com/rs/zerolog/log"
)

const dbSchema = `
CREATE TABLE IF NOT EXISTS CommandTikTokCustomUserVoices (
	person INTEGER NOT NULL,
	place INTEGER NOT NULL,
	voice VARCHAR(255) NOT NULL,
	UNIQUE(person, place)
)
`

func dbPersonVoiceExists(person, place int64) (bool, error) {
	db := core.DB
	db.Lock.Lock()
	defer db.Lock.Unlock()

	var exists bool

	row := db.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM CommandTikTokCustomUserVoices
			WHERE person = ? and place = ?
			LIMIT 1
		)`, person, place)

	err := row.Scan(&exists)

	log.Debug().
		Err(err).
		Int64("person", person).
		Int64("place", place).
		Bool("exists", exists).
		Msg("checked db to see if voice exists")

	return exists, err
}

func dbPersonGetVoice(person, place int64) (string, error) {
	db := core.DB
	db.Lock.RLock()
	defer db.Lock.RUnlock()

	var voice string

	row := db.DB.QueryRow(`
		SELECT voice
		FROM CommandTikTokCustomUserVoices
		WHERE person = ? and place = ?`, person, place)

	err := row.Scan(&voice)

	log.Debug().
		Err(err).
		Int64("person", person).
		Int64("place", place).
		Str("voice", voice).
		Msg("got voice for person")

	return voice, err
}

func dbPersonAddVoice(person, place int64, voice string) error {
	db := core.DB
	db.Lock.Lock()
	defer db.Lock.Unlock()

	_, err := db.DB.Exec(`
		INSERT INTO CommandTikTokCustomUserVoices(person, place, voice)
		VALUES (?, ?, ?)`, person, place, voice)

	log.Debug().
		Err(err).
		Int64("person", person).
		Int64("place", place).
		Str("voice", voice).
		Msg("added person voice in db")

	return err
}

func dbPersonUpdateVoice(person, place int64, voice string) error {
	db := core.DB
	db.Lock.Lock()
	defer db.Lock.Unlock()

	_, err := db.DB.Exec(`
		UPDATE CommandTikTokCustomUserVoices
		SET voice = ?
		WHERE person = ? and place = ?
	`, voice, person, place)

	log.Debug().
		Err(err).
		Int64("person", person).
		Int64("place", place).
		Str("voice", voice).
		Msg("updated person voice")

	return err
}
