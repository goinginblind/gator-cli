package cli

import (
	"fmt"
)

type Command struct {
	Name string
	Args []string
}

type commands struct {
	handlers map[string]func(*State, Command) error
}

func (c *commands) Run(s *State, cmd Command) error {
	f, ok := c.handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return f(s, cmd)
}
