package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadDB(v *viper.Viper, alias string) (DBConfig, error) {
	key := "databases." + alias

	if !v.IsSet(key) {
		return DBConfig{}, fmt.Errorf("database '%s' not found", alias)
	}

	var db DBConfig
	if err := v.UnmarshalKey(key, &db); err != nil {
		return DBConfig{}, fmt.Errorf("unmarshal config for database '%s': %w", alias, err)
	}

	return db, nil
}
