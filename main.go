package main

import (
	"log"
	"os"

	"github.com/rickNoise/aggreGATOR/app"
	"github.com/rickNoise/aggreGATOR/internal/config"
)

func main() {
	// Read the config file.
	c, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Store config in a new instance of the State struct.
	state := &app.State{Cfg: c}

	// Create a new instance of the commands struct with an initialized map of handler functions.
	commands := &app.Commands{
		RegisteredCommands: make(map[string]func(*app.State, app.Command) error),
	}

	// Register a handler function for the login command.
	commands.Register("login", app.HandlerLogin)

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
