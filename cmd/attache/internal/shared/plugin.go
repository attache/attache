package shared

type Command interface {
	Execute(flags []string) error
}

type Plugin interface {
	Name() string
	Command() Command
}

type pluginImpl struct {
	name string
	fn   func() Command
}

func NewPlugin(name string, fn func() Command) Plugin {
	return pluginImpl{name: name, fn: fn}
}

func (p pluginImpl) Name() string     { return p.name }
func (p pluginImpl) Command() Command { return p.fn() }
