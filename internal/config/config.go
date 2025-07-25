package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Represents the JSON config file structure, including struct tags.
type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	fmt.Printf("user name updated to %s\n", c.CurrentUserName)
	err := c.write()
	if err != nil {
		return fmt.Errorf("error writing config to file: %w", err)
	}
	return nil
}

func (c *Config) write() error {
	jsonData, err := json.MarshalIndent(*c, "", "\t")
	if err != nil {
		return fmt.Errorf("error marshalling json: %w", err)
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(configFilePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing json data to config file %s: %w", configFilePath, err)
	}

	return nil
}
