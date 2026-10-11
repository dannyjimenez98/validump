package config

type Config struct {
	Databases map[string]DBConfig `yaml:"databases"`
}

type DBConfig struct {
	Type string `yaml:"type"`
	ConnectionStr string `yaml:"connection_str"`
	BackupDir string `yaml:"backup_dir"`
}

