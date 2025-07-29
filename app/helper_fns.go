package app

import (
	"errors"
	"fmt"
)

/* HELPER FUNCTIONS */

// Helper function that updates the username in the config file.
func setUserNameInConfig(s *State, username string) error {
	if username == "" {
		return errors.New("cannot set a blank username")
	}
	err := s.Cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("could not set user name: %w", err)
	}
	return nil
}
