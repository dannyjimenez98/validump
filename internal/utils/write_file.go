package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// WriteAtomicFile creates a complete backup file without overwriting an existing file.
// Data is written to a temporary file at the dest before being moved to a permanent
// location at dest.
func WriteAtomicFile(dest string, write func(io.Writer) error) error {
	temp, err := os.CreateTemp(filepath.Dir(dest), "."+filepath.Base(dest)+".*.temp")
	if err != nil {
		return fmt.Errorf("create temporary backup location: %w", err)
	}
	closed := false
	defer func() {
		if !closed {
			temp.Close()
		}

		os.Remove(temp.Name())
	}()

	if err := write(temp); err != nil {
		return fmt.Errorf("write backup file: %w", err)
	}

	syncErr := temp.Sync()
	closeErr := temp.Close()
	closed = true

	if err := errors.Join(syncErr, closeErr); err != nil {
		return fmt.Errorf("sync backup file: %w", err)
	}

	// rename using Link instead of Rename because Link fails if new path already exists
	// not compatible with exFAT
	if err := os.Link(temp.Name(), dest); err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("destination already exists: %w", err)
		}

		return fmt.Errorf("commit backup: %w", err)
	}

	dir, err := os.Open(filepath.Dir(dest))
	if err != nil {
		return fmt.Errorf("backup committed, open parent directory: %w", err)
	}

	syncErr = dir.Sync()
	closeErr = dir.Close()

	if err := errors.Join(syncErr, closeErr); err != nil {
		return fmt.Errorf("backup committed, sync parent directory: %w", err)
	}

	return nil
}
