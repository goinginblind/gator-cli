package cli

import "maps"

var commandsMap = map[string]func(*state, Command) error{
	"login":    handlerLogin,
	"register": handlerRegister,
	"reset":    handlerReset,
	"users":    handlerGetUsers,
}

func NewCommands() *commands {
	cmds := &commands{handlers: make(map[string]func(*state, Command) error)}
	maps.Copy(cmds.handlers, commandsMap)
	return cmds
}
