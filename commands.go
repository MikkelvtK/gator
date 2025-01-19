package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if _, ok := c.commands[cmd.name]; !ok {
		return fmt.Errorf("commands: does not exist %s", cmd.name)
	}

	err := c.commands[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}
