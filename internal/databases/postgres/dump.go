package postgres

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Dump(w io.Writer, args []string) error {
	pgDumpPath, err := exec.LookPath("pg_dump")
	if err != nil {
		return fmt.Errorf("pg_dump not found in $PATH:" +
			"install PostgreSQL client tools or add their bin directory to PATH")
	}

	cmd := exec.Command(pgDumpPath, args...)
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
