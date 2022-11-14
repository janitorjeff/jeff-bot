package urban_dictionary

import (
	"git.slowtyper.com/slowtyper/janitorjeff/core"
	"git.slowtyper.com/slowtyper/janitorjeff/frontends"

	dg "github.com/bwmarrin/discordgo"
)

var Advanced = advanced{}

type advanced struct{}

func (advanced) Type() core.CommandType {
	return core.Advanced
}

func (advanced) Permitted(*core.Message) bool {
	return true
}

func (advanced) Names() []string {
	return []string{
		"ud",
	}
}

func (advanced) Description() string {
	return "Search a term or get a random one on urban dictionary."
}

func (c advanced) UsageArgs() string {
	return c.Children().Usage()
}

func (advanced) Parent() core.CommandStatic {
	return nil
}

func (advanced) Children() core.CommandsStatic {
	return core.CommandsStatic{
		AdvancedSearch,
		AdvancedRandom,
	}
}

func (advanced) Init() error {
	return nil
}

func (advanced) Run(m *core.Message) (any, error, error) {
	return m.Usage(), core.ErrMissingArgs, nil
}

////////////
//        //
// search //
//        //
////////////

var AdvancedSearch = advancedSearch{}

type advancedSearch struct{}

func (c advancedSearch) Type() core.CommandType {
	return c.Parent().Type()
}

func (c advancedSearch) Permitted(m *core.Message) bool {
	return c.Parent().Permitted(m)
}

func (advancedSearch) Names() []string {
	return []string{
		"search",
		"find",
	}
}

func (advancedSearch) Description() string {
	return "Search a term."
}

func (advancedSearch) UsageArgs() string {
	return "<term...>"
}

func (advancedSearch) Parent() core.CommandStatic {
	return Advanced
}

func (advancedSearch) Children() core.CommandsStatic {
	return nil
}

func (advancedSearch) Init() error {
	return nil
}

func (c advancedSearch) Run(m *core.Message) (any, error, error) {
	if len(m.Command.Args) < 1 {
		return m.Usage(), core.ErrMissingArgs, nil
	}

	switch m.Frontend {
	case frontends.Discord:
		return c.discord(m)
	default:
		return c.text(m)
	}
}

func (c advancedSearch) discord(m *core.Message) (*dg.MessageEmbed, error, error) {
	def, usrErr, err := c.core(m)
	if err != nil {
		return nil, nil, err
	}
	return renderDiscord(def, usrErr), usrErr, nil
}

func (c advancedSearch) text(m *core.Message) (string, error, error) {
	def, usrErr, err := c.core(m)
	if err != nil {
		return "", nil, err
	}
	return renderText(def, usrErr), usrErr, nil
}

func (advancedSearch) core(m *core.Message) (definition, error, error) {
	term := m.RawArgs(0)
	return search(term)
}

////////////
//        //
// random //
//        //
////////////

var AdvancedRandom = advancedRandom{}

type advancedRandom struct{}

func (c advancedRandom) Type() core.CommandType {
	return c.Parent().Type()
}

func (c advancedRandom) Permitted(m *core.Message) bool {
	return c.Parent().Permitted(m)
}

func (advancedRandom) Names() []string {
	return []string{
		"random",
		"rand",
	}
}

func (advancedRandom) Description() string {
	return "Get a random term."
}

func (advancedRandom) UsageArgs() string {
	return ""
}

func (advancedRandom) Parent() core.CommandStatic {
	return Advanced
}

func (advancedRandom) Children() core.CommandsStatic {
	return nil
}

func (advancedRandom) Init() error {
	return nil
}

func (c advancedRandom) Run(m *core.Message) (any, error, error) {
	switch m.Frontend {
	case frontends.Discord:
		return c.discord()
	default:
		return c.text()
	}
}

func (c advancedRandom) discord() (*dg.MessageEmbed, error, error) {
	def, usrErr, err := c.core()
	if err != nil {
		return nil, nil, err
	}
	return renderDiscord(def, usrErr), usrErr, nil
}

func (c advancedRandom) text() (string, error, error) {
	def, usrErr, err := c.core()
	if err != nil {
		return "", nil, err
	}
	return renderText(def, usrErr), usrErr, nil
}

func (advancedRandom) core() (definition, error, error) {
	return rand()
}
