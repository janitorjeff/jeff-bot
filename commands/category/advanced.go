package category

import (
	"fmt"

	"github.com/janitorjeff/jeff-bot/core"
	"github.com/janitorjeff/jeff-bot/frontends/twitch"
)

var Advanced = advanced{}

type advanced struct{}

func (advanced) Type() core.CommandType {
	return core.Advanced
}

func (advanced) Permitted(m *core.Message) bool {
	if m.Frontend.Type() != twitch.Frontend.Type() {
		return false
	}
	return m.Author.Mod()
}

func (advanced) Names() []string {
	return []string{
		"category",
		"game",
	}
}

func (advanced) Description() string {
	return "Show or edit the current category."
}

func (c advanced) UsageArgs() string {
	return c.Children().Usage()
}

func (advanced) Category() core.CommandCategory {
	return core.CommandCategoryModerators
}

func (advanced) Examples() []string {
	return nil
}

func (advanced) Parent() core.CommandStatic {
	return nil
}

func (advanced) Children() core.CommandsStatic {
	return core.CommandsStatic{
		AdvancedShow,
		AdvancedEdit,
	}
}

func (advanced) Init() error {
	return nil
}

func (advanced) Run(m *core.Message) (any, error, error) {
	return m.Usage(), core.ErrMissingArgs, nil
}

//////////
//      //
// show //
//      //
//////////

var AdvancedShow = advancedShow{}

type advancedShow struct{}

func (c advancedShow) Type() core.CommandType {
	return c.Parent().Type()
}

func (c advancedShow) Permitted(m *core.Message) bool {
	return c.Parent().Permitted(m)
}

func (advancedShow) Names() []string {
	return core.AliasesShow
}

func (advancedShow) Description() string {
	return "Show the current category."
}

func (advancedShow) UsageArgs() string {
	return ""
}

func (c advancedShow) Category() core.CommandCategory {
	return c.Parent().Category()
}

func (c advancedShow) Examples() []string {
	return nil
}

func (advancedShow) Parent() core.CommandStatic {
	return Advanced
}

func (advancedShow) Children() core.CommandsStatic {
	return nil
}

func (advancedShow) Init() error {
	return nil
}

func (c advancedShow) Run(m *core.Message) (any, error, error) {
	switch m.Frontend.Type() {
	case twitch.Frontend.Type():
		return c.twitch(m)
	default:
		panic("this should never happen")
	}
}

func (advancedShow) twitch(m *core.Message) (string, error, error) {
	h, err := m.Client.(*twitch.Twitch).Helix()
	if err != nil {
		return "", nil, err
	}

	g, err := h.GetGameName(m.Here.ID())
	return g, nil, err
}

//////////
//      //
// edit //
//      //
//////////

var AdvancedEdit = advancedEdit{}

type advancedEdit struct{}

func (c advancedEdit) Type() core.CommandType {
	return c.Parent().Type()
}

func (c advancedEdit) Permitted(m *core.Message) bool {
	return c.Parent().Permitted(m)
}

func (advancedEdit) Names() []string {
	return core.AliasesEdit
}

func (advancedEdit) Description() string {
	return "Edit the current category."
}

func (advancedEdit) UsageArgs() string {
	return "<category...>"
}

func (c advancedEdit) Category() core.CommandCategory {
	return c.Parent().Category()
}

func (c advancedEdit) Examples() []string {
	return []string{
		"minecraft",
		"just chatting",
	}
}

func (advancedEdit) Parent() core.CommandStatic {
	return Advanced
}

func (advancedEdit) Children() core.CommandsStatic {
	return nil
}

func (advancedEdit) Init() error {
	return nil
}

func (c advancedEdit) Run(m *core.Message) (any, error, error) {
	switch m.Frontend.Type() {
	case twitch.Frontend.Type():
		return c.twitch(m)
	default:
		panic("this should never happen")
	}
}

func (advancedEdit) twitch(m *core.Message) (string, error, error) {
	h, err := m.Client.(*twitch.Twitch).Helix()
	if err != nil {
		return "", nil, err
	}

	g, usrErr, err := h.SetGame(m.Here.ID(), m.RawArgs(0))

	if usrErr != nil {
		return fmt.Sprint(usrErr), usrErr, nil
	}

	switch err {
	case nil:
		return fmt.Sprintf("Category set to: %s", g), nil, nil
	case twitch.ErrNoResults:
		return "Couldn't find the category, did you type the name correctly?", nil, nil
	default:
		return "", nil, err
	}
}
