package app

import (
	"context"
	"fmt"

	"github.com/rickNoise/aggreGATOR/internal/database"
)

// Create logged-in middleware.
// It will allow us to change the function signature of our handlers that require a logged in user to accept a user as an argument and DRY up our code.
// You'll notice it's a higher order function that takes a handler of the "logged in" type and returns a "normal" handler that we can register. I used it like this:
// cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("could not get current user, maybe no user is logged in?: %w", err)
		}
		return handler(s, cmd, user)
	}
}
