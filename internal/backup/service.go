package backup

import "github.com/dannyjimenez98/validump/internal/config"

type Options struct {
	Format      string
	Compression bool
}

type Service struct {
	alias     string
	db        config.DBConfig
	backupDir string
	opts      Options
}

func NewService(alias string, cfg config.DBConfig, backupDir string, opts Options) *Service {
	return &Service{
		alias:     alias,
		db:        cfg,
		backupDir: backupDir,
		opts:      opts,
	}
}
