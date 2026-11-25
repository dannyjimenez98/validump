package utils

import (
	"compress/gzip"
	"errors"
	"io"
)

func Compress(dest string, write func(io.Writer) error) error {
	return WriteAtomicFile(dest, func(w io.Writer) error {
		zw := gzip.NewWriter(w)
		writeErr := write(zw)
		closeErr := zw.Close()
		return errors.Join(writeErr, closeErr)
	})
}
