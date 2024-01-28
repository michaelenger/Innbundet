package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Configuration for the innbundet server.
type Config struct {
	Title        string
	Description  string
	DatabaseFile string `yaml:"database_file"`
	ItemsPerPage int    `yaml:"items_per_page"`
}

// Read the config from a file path.
func FromFile(filePath string) (*Config, error) {
	logger := log.Default()

	conf := Config{
		"Innbundet",
		"Tiny RSS reader.",
		"innbundet.sqlite",
		25,
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logger.Printf("Unable to read config file at: %s; using default values", filePath)
		return &conf, nil
	}

	logger.Printf("Reading config from: %s", filePath)
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(contents), &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
