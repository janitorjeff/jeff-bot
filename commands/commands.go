package commands

import (
	"git.slowtyper.com/slowtyper/janitorjeff/commands/command"
	"git.slowtyper.com/slowtyper/janitorjeff/commands/help"
	"git.slowtyper.com/slowtyper/janitorjeff/commands/nick"
	"git.slowtyper.com/slowtyper/janitorjeff/commands/paintball"
	"git.slowtyper.com/slowtyper/janitorjeff/commands/prefix"
	"git.slowtyper.com/slowtyper/janitorjeff/commands/rps"
	"git.slowtyper.com/slowtyper/janitorjeff/commands/scope"
	"git.slowtyper.com/slowtyper/janitorjeff/commands/test"
	"git.slowtyper.com/slowtyper/janitorjeff/commands/time"
	"git.slowtyper.com/slowtyper/janitorjeff/commands/urban-dictionary"
	"git.slowtyper.com/slowtyper/janitorjeff/commands/wikipedia"
	"git.slowtyper.com/slowtyper/janitorjeff/core"
)

var Normal = core.Commands{
	command.Command,
	help.Normal,
	nick.Normal,
	paintball.Normal,
	prefix.Command,
	rps.Normal,
	time.Normal,
	urban_dictionary.Normal,
	wikipedia.Normal,
}

var Advanced = core.Commands{
	help.Advanced,
	nick.Advanced,
	time.Advanced,
	urban_dictionary.Advanced,
}

var Admin = core.Commands{
	help.Admin,
	nick.Admin,
	prefix.Admin,
	scope.Admin,
	test.Command,
}
