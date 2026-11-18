package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

func (svc *Service) Backup() error {
	if svc.db.Engine == "postgres" {
		timestamp := time.Now().UTC().Format("2006-01-02T15-04-05Z")
		backupFile := fmt.Sprintf("%s_%s.dump", svc.alias, timestamp)

		pgDumpPath, err := exec.LookPath("pg_dump")
		if err != nil {
			return fmt.Errorf(
				"backup database '%s': pg_dump not found in $PATH: "+
					"install PostgreSQL client tools or add their bin directory to PATH",
				svc.alias,
			)
		}

		cmd := exec.Command(
			pgDumpPath,
			"-h", svc.db.Host,
			"-p", strconv.Itoa(int(svc.db.Port)),
			"-U", svc.db.User,
			"-d", svc.db.DBName,
			"-F", svc.opts.Format,
			"-Z", svc.opts.Compression,
			"-f", filepath.Join(svc.outputPath, backupFile),
		)

		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("run backup command: %w", err)
		}

		fmt.Println("Backup completed successfully!")

		return nil
	}

	return fmt.Errorf("unsupported database: %s ('%s')", svc.db.Engine, svc.alias)
}
