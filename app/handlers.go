package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/rickNoise/aggreGATOR/internal/database"
	"github.com/rickNoise/aggreGATOR/rss"
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

// HandlerReset deletes all data in the users DB table.
func HandlerReset(s *State, cmd Command) error {
	if len(cmd.Arguments) != 0 {
		return errors.New("reset failed: reset cannot pass any arguments")
	}

	err := s.Db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Successfully reset database; all users data deleted")
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	if len(cmd.Arguments) != 0 {
		return errors.New("users command cannot take arguments")
	}

	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("no users registered")
		return nil
	}

	for _, user := range users {
		if user.Name == s.Cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	// hardcoded single URL for now
	const URL = "https://www.wagslane.dev/index.xml"

	feed, err := rss.FetchFeed(context.Background(), URL)
	if err != nil {
		return err
	}

	fmt.Printf("successfully fetched feed from %s: %+v\n", URL, *feed)
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Arguments) < 2 {
		return errors.New("error: addfeed command must pass two arguments (name) and (url)")
	}
	feedName := cmd.Arguments[0]
	feedUrl := cmd.Arguments[1]

	user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not get current user: %w", err)
	}

	feed, err := s.Db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      feedName,
			Url:       feedUrl,
			UserID:    user.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("could not create feed in the database: %w", err)
	}

	fmt.Printf("successfully added feed %s to user %s: %+v", feed.Name, user.Name, feed)
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	if len(cmd.Arguments) != 0 {
		return errors.New("feeds command cannot take arguments")
	}

	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds data from db: %w", err)
	}
	if len(feeds) == 0 {
		fmt.Println("No feeds in database.")
	}

	/*
		The name of the feed
		The URL of the feed
		The name of the user that created the feed (you might need a new SQL query)
	*/
	for _, feed := range feeds {
		feedAdder, err := s.Db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("could not get creating user for feed %v, may be problem with feed data or GetUserById method: %w", feed, err)
		}

		fmt.Println("name:", feed.Name)
		fmt.Printf("├── url:      %s\n", feed.Url)
		fmt.Printf("└── added by: %s\n", feedAdder.Name)
	}
	return nil
}

// Takes a single url argument and creates a new feed follow record for the current user. It should print the name of the feed and the current user once the record is created (which the query we just made should support). You'll need a query to look up feeds by URL.
func HandlerFollow(s *State, cmd Command) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("follow command takes a single url argument")
	}
	url := cmd.Arguments[0]

	user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not get current user: %w", err)
	}
	feed, err := s.Db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not GetFeedByUrl, feed may not exist yet, try to add the feed first: %w", err)
	}

	feed_follow, err := s.Db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("error with CreateFeedFollow: %w", err)
	}

	fmt.Printf("Current user %s now following %s\n", feed_follow.UserName, feed_follow.FeedName)
	return nil
}
