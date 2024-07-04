package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Team is a list of e-mails for every member of the team.
type Team []string

type AuthorFilter struct {
	// ExcludeName is a regular expression dictating which author names should be excluded.
	ExcludeName string `yaml:"exclude_name"`
	// ExcludeEmail is a regular expression dictating which author emails should be excluded.
	ExcludeEmail string `yaml:"exclude_email"`
}

type Config struct {
	// AuthorFilters are filters that will remove the specified authors from the game.
	AuthorFilters []AuthorFilter `yaml:"author_filters"`
	// Teams is a map from team name to the members of that team.
	Teams map[string]Team `yaml:"teams"`
	// AuthorBias is how much to bias towards authors with high commit counts.
	AuthorBias *float64 `yaml:"author_bias"`
}

func Load() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	path := filepath.Join(home, ".gauthordle.yaml")
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file simply doesn't exist, we provide an empty config.
			return Config{}, nil
		}
		return Config{}, err
	}
	defer f.Close()

	var cfg Config
	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
