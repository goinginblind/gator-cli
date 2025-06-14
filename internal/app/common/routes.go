package common

import (
	"fmt"
)

type Command struct {
	Name string
	Args []string
}

type Routes struct {
	handlers map[string]func(*State, Command) error
}

func NewRoutes() *Routes {
	return &Routes{handlers: make(map[string]func(*State, Command) error)}
}

func (c *Routes) Run(s *State, cmd Command) error {
	f, ok := c.handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return f(s, cmd)
}

func (c *Routes) Register(cmd string, f func(*State, Command) error) {
	c.handlers[cmd] = f
}
