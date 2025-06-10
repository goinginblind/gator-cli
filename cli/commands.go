package cli

import (
	"fmt"
)

type Command struct {
	Name string
	Args []string
}

type commands struct {
	handlers map[string]func(*state, Command) error
}

func (c *commands) Run(s *state, cmd Command) error {
	f, ok := c.handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, Command) error) {
	c.handlers[name] = f
}
