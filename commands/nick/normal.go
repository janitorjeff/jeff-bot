package nick

import (
	"git.slowtyper.com/slowtyper/janitorjeff/core"
)

var Normal = &core.CommandStatic{
	Names: []string{
		"nick",
		"nickname",
	},
	Description: "View or set your nickname.",
	UsageArgs:   "[nickname]",
	Run:         normalRun,
}

func normalRun(m *core.Message) (any, error, error) {
	if len(m.Command.Runtime.Args) == 0 {
		return advancedRunView(m)
	}
	return advancedRunSet(m)
}

// Tries to find a user scope from the given string. First tries to find if it
// matches a nickname in the database and if it doesn't it tries various
// platform specific things, like for example checking if the given string is a
// user ID.
func ParseUser(m *core.Message, place int64, s string) (int64, error) {
	if user, err := dbGetUser(s, place); err == nil {
		return user, nil
	}

	placeID, err := core.Globals.DB.ScopeID(place)
	if err != nil {
		return -1, err
	}

	id, err := m.Client.PersonID(s, placeID)
	if err != nil {
		return -1, err
	}

	return m.Client.PersonScope(id)
}

// Same as ParseUser but uses the default place instead
func ParseUserHere(m *core.Message, s string) (int64, error) {
	place, err := m.ScopeHere()
	if err != nil {
		return -1, err
	}

	return ParseUser(m, place, s)
}
