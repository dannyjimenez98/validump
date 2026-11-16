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

		if err := validateDBConfig(dbConfig); err != nil {
			return fmt.Errorf("invalid configuration for '%s': %w", dbAlias, err)
		}
	}
	return nil
}

func validateDBConfig(dbConfig map[string]any) error {
	engine, ok := dbConfig["engine"].(string)
	if !ok || strings.TrimSpace(engine) == "" {
		return fmt.Errorf("'engine' must be a non-empty string")
	}

	var requiredFields []string
	var requiresPort bool

	switch engine {
	case "postgres":
		requiredFields = []string{"host", "username", "dbname"}
		requiresPort = true

	default:
		return fmt.Errorf("unsupported database engine %q", engine)
	}

	for _, field := range requiredFields {
		value, exists := dbConfig[field]
		if !exists {
			return fmt.Errorf("required field %q is missing", field)
		}
		text, ok := value.(string)
		if !ok || strings.TrimSpace(text) == "" {
			return fmt.Errorf("field %q must be a non-empty string", field)
		}
	}

	// validate port when required or explicitly provided
	portValue, hasPort := dbConfig["port"]
	if requiresPort && !hasPort {
		return fmt.Errorf("required field 'port' is missing")
	}
	if hasPort {
		port, ok := portValue.(int)
		if !ok || port < 1 || port > 65535 {
			return fmt.Errorf("'port' must be an integer between 1 and 65535")
		}
	}

	// passwords are optional because they may come from the environment
	if password, exists := dbConfig["password"]; exists {
		if _, ok := password.(string); !ok {
			return fmt.Errorf("'password' must be a string")
		}
	}

	return nil
}
