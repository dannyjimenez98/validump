package config

type Config struct {
	BackupDir string              `mapstructure:"backup_dir"`
	Databases map[string]DBConfig `mapstructure:"databases"`
}

type DBConfig struct {
	Engine   string `mapstructure:"engine"`
	Host     string `mapstructure:"host"`
	Port     uint16 `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}
