package app

import (
	"errors"
	"fmt"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return errors.New("login expects a username argument")
	}

	// Set the user to the given username
	err := s.Cfg.SetUser(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("could not set user name: %w", err)
	}
	fmt.Printf("user has been set to %s\n", cmd.Arguments[0])
	return nil
}
