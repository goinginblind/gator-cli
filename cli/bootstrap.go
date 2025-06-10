package cli

var commandsMap = map[string]func(*state, Command) error{
	"login": handlerLogin,
}

func NewCommands() *commands {
	cmds := &commands{handlers: make(map[string]func(*state, Command) error)}
	for k, v := range commandsMap {
		cmds.handlers[k] = v
	}
	return cmds
}
