package main

// test add, can be removed

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/rickNoise/aggreGATOR/app"
	"github.com/rickNoise/aggreGATOR/internal/config"
	"github.com/rickNoise/aggreGATOR/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	// Read the config file.
	c, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Load in your database URL to the config struct and sql.Open() a connection to your database.
	db, err := sql.Open("postgres", c.DbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Use your generated database package to create a new *database.Queries, and store it in your state struct.
	dbQueries := database.New(db)

	// Store config in a new instance of the State struct.
	state := &app.State{
		Db:  dbQueries,
		Cfg: c,
	}

	// Create a new instance of the commands struct with an initialized map of handler functions.
	commands := &app.Commands{
		RegisteredCommands: make(map[string]func(*app.State, app.Command) error),
	}

	// Register handles for cli commands
	commands.Register("login", app.HandlerLogin)
	commands.Register("register", app.HandlerRegister)
	commands.Register("reset", app.HandlerReset)
	commands.Register("users", app.HandlerUsers)
	commands.Register("agg", app.HandlerAgg)
	commands.Register("addfeed", app.MiddlewareLoggedIn(app.HandlerAddFeed))
	commands.Register("feeds", app.HandlerFeeds)
	commands.Register("follow", app.MiddlewareLoggedIn(app.HandlerFollow))
	commands.Register("following", app.MiddlewareLoggedIn(app.HandlerFollowing))
	commands.Register("unfollow", app.MiddlewareLoggedIn(app.HandlerUnfollow))
	commands.Register("browse", app.MiddlewareLoggedIn(app.HandlerBrowse))

	// Register handle for help command
	commands.Register("help", func(s *app.State, cmd app.Command) error {
		for name := range commands.RegisteredCommands {
			fmt.Println(name)
		}
		return nil
	})

	// Use os.Args to get the command-line arguments passed in by the user.
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	userCommand := app.Command{
		Name:      os.Args[1],
		Arguments: os.Args[2:],
	}

	err = commands.Run(state, userCommand)
	if err != nil {
		log.Fatalf("could not run command: %v", err)
	}
}
