package prefix

import (
	"git.slowtyper.com/slowtyper/janitorjeff/core"
)

type flags struct {
	fs *core.Flags

	typeFlag  core.Type
	scopeFlag int64
}

func newFlags(m *core.Message) *flags {
	f := &flags{
		fs: core.NewFlags(m),
	}
	return f
}

func (f *flags) TypeFlag() *flags {
	core.TypeFlag(&f.typeFlag, core.All, f.fs)
	return f
}

func (f *flags) ScopeFlag() *flags {
	core.FlagPlace(&f.scopeFlag, f.fs)
	return f
}
