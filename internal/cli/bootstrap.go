package cli

import "maps"

var commandsMap = map[string]func(*State, Command) error{
	"login":     handlerLogin,
	"register":  handlerRegister,
	"reset":     handlerReset,
	"users":     handlerGetUsers,
	"agg":       handlerAggregator,
	"addfeed":   handlerCreateFeed,
	"feeds":     handlerGetFeedsWithUNames,
	"follow":    handlerCreateFollow,
	"following": handlerGetFeedFollowsForUser,
}

func NewCommands() *commands {
	cmds := &commands{handlers: make(map[string]func(*State, Command) error)}
	maps.Copy(cmds.handlers, commandsMap)
	return cmds
}
