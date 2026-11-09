package config

type Config struct {
	BackupDir string              `yaml:"backup_dir"`
	Databases map[string]DBConfig `yaml:"databases"`
}

type DBConfig struct {
	Engine   string `yaml:"engine"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	DBName   string `yaml:"dbname"`
}
