package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/rickNoise/aggreGATOR/internal/database"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return errors.New("login expects a username argument")
	}
	providedUsername := cmd.Arguments[0]

	_, err := s.Db.GetUser(context.Background(), providedUsername)
	if err != nil {
		return fmt.Errorf("could not GetUser, user likely does not exist: %w", err)
	}

	// Set the user to the given username
	err = setUserNameInConfig(s, providedUsername)
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return errors.New("register expects a username argument")
	}
	providedUsername := cmd.Arguments[0]

	user, err := s.Db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      providedUsername,
		},
	)
	if err != nil {
		log.Fatalf("error creating user, likely user already exists: %v", err)
	}

	err = setUserNameInConfig(s, user.Name)
	if err != nil {
		return fmt.Errorf("user registered, but could not be set in config: %w", err)
	}

	// Print a message that the user was created, and log the user's data to the console for your own debugging.
	fmt.Printf("user successfuly registered: %+v", user)
	return nil
}

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
