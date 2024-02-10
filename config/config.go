package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
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
	conf := Config{
		"Innbundet",
		"Tiny RSS reader.",
		"innbundet.sqlite",
		25,
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Warn().
			Str("path", filePath).
			Msg("Unable to find config file; using default values")
		return &conf, nil
	}
	log.Info().
		Str("path", filePath).
		Msg("Reading config file")

	contents, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to read config file: %s", err))
	}

	err = yaml.Unmarshal([]byte(contents), &conf)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to unmarshal config file: %s", err))
	}

	return &conf, nil
}
