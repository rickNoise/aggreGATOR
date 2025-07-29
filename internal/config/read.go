package config

import (
	"encoding/json"
	"fmt"
	"os"
)

/* CONSTANTS */
const configFileName = ".gatorconfig.json"

// Reads the JSON config file found at jsonFilePath and returns its contents a Config struct.
func Read() (*Config, error) {
	// Set the filepath to the user's home directory + filename constant at the top of this file
	jsonFilePath, err := getConfigFilePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config file path: %w", err)
	}

	// Read the JSON file into a byte slice directly
	// fmt.Printf("Reading config from %s...\n", jsonFilePath)
	fileData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", jsonFilePath, err)
	}

	// Unmarshal the JSON data into our Go struct
	var config Config
	err = json.Unmarshal(fileData, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from %s: %w", jsonFilePath, err)
	}

	return &config, nil
}

func getConfigFilePath() (string, error) {
	// from a constant at top of this file
	fileName := configFileName

	if homeDir, err := os.UserHomeDir(); err != nil {
		return "", fmt.Errorf("failed to get user home dir: %w", err)
	} else {
		return homeDir + "/" + fileName, nil
	}
}
