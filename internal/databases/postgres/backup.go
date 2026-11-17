package postgres

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dannyjimenez98/validump/internal/config"
)

func Backup(cfg config.DBConfig, path string, dbalias string) error {
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05Z")
	backupFile := fmt.Sprintf("%s_%s.dump", dbalias, timestamp)

	pgDumpPath, err := exec.LookPath("pg_dump")
	if err != nil {
		return fmt.Errorf("backup database '%s': pg_dump not found in $PATH: install PostgreSQL client tools or add their bin directory to PATH", dbalias)
	}

	cmd := exec.Command(
		pgDumpPath,
		"-h", cfg.Host,
		"-p", strconv.Itoa(int(cfg.Port)),
		"-U", cfg.Username,
		"-d", cfg.DBName,
		"-Fc",
		"-f", filepath.Join(path, backupFile),
	)

	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run backup command: %w", err)
	}

	fmt.Println("Backup completed successfully!")

	return nil
}
