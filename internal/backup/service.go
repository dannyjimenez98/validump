package backup

import "github.com/dannyjimenez98/validump/internal/config"

type Options struct {
	Format      string
	Compression string
}

type Service struct {
	alias      string
	db         config.DBConfig
	outputPath string
	opts       Options
}

func NewService(alias string, cfg config.DBConfig, outputPath string, opts Options) *Service {
	return &Service{
		alias:      alias,
		db:         cfg,
		outputPath: outputPath,
		opts:       opts,
	}
}
