package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func Load(configPath string) (*viper.Viper, error) {
	v := viper.New()
	v.SetEnvPrefix("VALIDUMP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("file not found at path %s. Run validump init [--config <path>]: %w", configPath, err)
		} else {
			return nil, fmt.Errorf("read configuration file: %w", err)
		}
	}

	fmt.Printf("Using config file: %s\n", v.ConfigFileUsed())

	return v, nil
}
