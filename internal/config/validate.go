package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v4"
)

func Validate(cfgFile string) error {
	yamlData, err := os.ReadFile(cfgFile)
	if err != nil {
		return fmt.Errorf("read YAML config file: %w", err)
	}

	var doc any
	if err := yaml.Unmarshal(yamlData, &doc); err != nil {
		return fmt.Errorf("unmarshal YAML data: %w", err)
	}

	if doc == nil {
		return fmt.Errorf("configuration cannot be empty or null")
	}

	root, ok := doc.(map[string]any)
	if !ok {
		return fmt.Errorf("configuration does not contain YAML mapping")
	}

	val, exists := root["databases"]
	if !exists {
		return fmt.Errorf("required 'databases' key could not be found in the configuration")
	}

	databases, ok := val.(map[string]any)
	if !ok {
		return fmt.Errorf("'databases' must be a YAML mapping")
	}

	if len(databases) == 0 {
		return fmt.Errorf("'databases' must contain at least one child database mapping")
	}

	// track duplicate aliases
	seen := make(map[string]bool)

	// checks aliases against valid formatting:
	pattern := "^[a-z](?:[a-z0-9]|_[a-z0-9])*$"
	regex := regexp.MustCompile(pattern)

	for dbAlias, val := range databases {
		if len(dbAlias) > 50 {
			return fmt.Errorf("database alias '%s' exceeds 50 characters", dbAlias)
		}

		normalized := strings.ToLower(dbAlias)
		if seen[normalized] {
			return fmt.Errorf("duplicate database alias: '%s'", dbAlias)
		}
		seen[normalized] = true

		if !regex.MatchString(dbAlias) {
			return fmt.Errorf("invalid database alias: '%s' must start with a lowercase letter, "+
				"contain only lowercase letters, numbers, or single underscores, "+
				"and cannot end with an underscore", dbAlias)
		}

		dbConfig, ok := val.(map[string]any)
		if !ok {
			return fmt.Errorf("'%s' configuration must be a mapping", dbAlias)
		}

		if len(dbConfig) == 0 {
			return fmt.Errorf("database '%s' must contain a configuration mapping", dbAlias)
		}
	}
	return nil
}
