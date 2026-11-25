package postgres

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Dump(ctx context.Context, args []string, w io.Writer) error {
	pgDumpPath, err := exec.LookPath("pg_dump")
	if err != nil {
		return fmt.Errorf("pg_dump not found in $PATH:" +
			"install PostgreSQL client tools or add their bin directory to PATH")
	}

	// cmd := exec.Command(pgDumpPath, args...)
	cmd := exec.CommandContext(ctx, pgDumpPath, args...)
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("backup interrupted: %w", errors.Join(err, ctx.Err()))
		}
		return fmt.Errorf("backup failed: %w", err)
	}

	return nil
}
