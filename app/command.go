package app

import (
	"errors"
)

// A Command contains a name and a slice of string arguments.
// For example, in the case of the login Command, the name would be "login" and the handler will expect the arguments slice to contain one string, the username.
type Command struct {
	Name      string
	Arguments []string
}

// This will hold all the Commands the CLI can handle.
type Commands struct {
	// Map of command names to their handler functions.
	CommandToHandlerMap map[string]func(*State, Command) error
}

// Runs a given Command with the provided state if it exists.
func (c *Commands) Run(s *State, cmd Command) error {
	fn, exists := c.CommandToHandlerMap[cmd.Name]
	if !exists {
		return errors.New("command does not exist")
	}

	err := fn(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

// Registers a new handler function for a command name.
func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CommandToHandlerMap[name] = f
}
