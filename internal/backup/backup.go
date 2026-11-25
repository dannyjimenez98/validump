package backup

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dannyjimenez98/validump/internal/databases/postgres"
	"github.com/dannyjimenez98/validump/internal/utils"
)

func (svc *Service) Backup() error {
	dbalias := svc.alias
	db := svc.db
	opts := svc.opts
	backupDir := svc.backupDir

	if db.Engine == "postgres" {
		filename, err := generateFileName(dbalias, opts.Format)
		if err != nil {
			return fmt.Errorf("generate filename: %w", err)
		}

		backupPath := filepath.Join(backupDir, filename)

		args := []string{
			"-h", db.Host,
			"-p", strconv.Itoa(int(db.Port)),
			"-U", db.User,
			"-d", db.DBName,
			"-F", opts.Format,

			// disable postgres default compression
			// if compression is requested, Validump's compression will be used
			"-Z", "0",
		}

		dump := func(w io.Writer) error {
			return postgres.Dump(w, args)
		}

		if opts.Format == "directory" {
			tempDir, err := os.MkdirTemp(backupDir, ".tempBackup-*_"+filename)
			if err != nil {
				return err
			}
			cleanupTemp := true
			defer func() {
				if cleanupTemp {
					os.RemoveAll(tempDir)
				}
			}()

			args = append(args, "-f", tempDir)

			if err := dump(io.Discard); err != nil {
				return err
			}

			if opts.Compression {
				err := utils.Compress(
					backupPath+".tar.gz",
					func(w io.Writer) error {
						return utils.Archive(tempDir, w)
					},
				)
				if err != nil {
					// block deferred cleanup to preserve temp backup
					cleanupTemp = false

					// preserve temp directory backup source if compression fails
					recoveryPath := backupPath + ".uncompressed"
					if preserveErr := os.Rename(tempDir, recoveryPath); preserveErr != nil {
						return fmt.Errorf("compress backup: %w: preserve backup source: %w - source retained at %s", err, preserveErr, tempDir)
					}

					return fmt.Errorf(
						"compress backup: %w"+" - source retained at %s",
						err,
						recoveryPath,
					)
				}
			} else {
				if err := os.Rename(tempDir, backupPath); err != nil {
					cleanupTemp = false
					return fmt.Errorf("commit directory backup: %w"+" - source retained at %s", err, tempDir)
				}
			}

		} else {
			write := utils.WriteAtomicFile
			if opts.Compression {
				backupPath += ".gz"
				write = utils.Compress
			}

			if err := write(backupPath, dump); err != nil {
				return fmt.Errorf("backup database '%s': %w", dbalias, err)
			}
		}

		fmt.Println("Backup completed successfully!")

		return nil
	}

	return fmt.Errorf("unsupported database: %s ('%s')", db.Engine, dbalias)
}

func generateFileName(alias string, format string) (string, error) {
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05Z")

	extensions := map[string]string{
		"plain":     ".sql",
		"custom":    ".dump",
		"directory": "",
		"tar":       ".tar",
	}

	ext, ok := extensions[format]
	if !ok {
		return "", fmt.Errorf("unsupported backup format: %s", format)
	}

	// append a 6 character hex string to differentiate backups created in same second
	var suffix [3]byte
	_, _ = rand.Read(suffix[:])

	filename := fmt.Sprintf("%s_%s_%s%s", alias, timestamp, hex.EncodeToString(suffix[:]), ext)

	return filename, nil
}
