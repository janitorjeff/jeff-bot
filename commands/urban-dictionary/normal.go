package urban_dictionary

import (
	"git.slowtyper.com/slowtyper/janitorjeff/core"
)

var Normal = normal{}

type normal struct{}

func (normal) Type() core.CommandType {
	return core.Normal
}

func (normal) Permitted(*core.Message) bool {
	return true
}

func (normal) Names() []string {
	return Advanced.Names()
}

func (normal) Description() string {
	return "Search a term on urban dictionary."
}

func (normal) UsageArgs() string {
	return "<term...>"
}

func (normal) Parent() core.CommandStatic {
	return nil
}

func (normal) Children() core.CommandsStatic {
	return nil
}

func (normal) Init() error {
	return nil
}

func (normal) Run(m *core.Message) (any, error, error) {
	return AdvancedSearch.Run(m)
}
