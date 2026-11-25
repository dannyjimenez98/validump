package utils

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Archive takes the directory at srcDir and its contents and writes it to a single tar file.
func Archive(srcDir string, w io.Writer) error {
	tw := tar.NewWriter(w)

	walkErr := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == srcDir {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if !info.IsDir() && !info.Mode().IsRegular() {
			return fmt.Errorf("unsupported file: %s", path)
		}

		// create header from file metadata
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		// store the header name as the relative path
		// and normalize slashes for cross-platform compatibility
		relativePath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relativePath)

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// call d.isDir after metadata is archived but return before
		// file contents are copied since its not a file
		if d.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := io.Copy(tw, file); err != nil {
			return err
		}

		return nil
	})

	return errors.Join(walkErr, tw.Close())
}
