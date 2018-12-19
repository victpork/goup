// +build !windows

package goup

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	Format            = "tar.gz"
	DefaultInstallDir = "/usr/local/go/bin"
)

func ExtractArchive(srcFile *os.File, size int64, targetPath string, progCback func(format string, arg ...interface{})) error {
	_, err := srcFile.Seek(0, 0)
	if err != nil {
		return errors.Wrap(err, "Error resetting offset")
	}
	gzFile, err := gzip.NewReader(srcFile)
	if err != nil {
		return errors.Wrap(err, "Error opening GZip archive")
	}
	tarFile := tar.NewReader(gzFile)
	for {
		f, err := tarFile.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "Error opening tar file")
		}
		dstPath := filepath.Join(targetPath, f.Name[3:])
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(filepath.Dir(dstPath), f.FileInfo().Mode())
			if err != nil {
				return errors.Wrap(err, "Cannot create directory")
			}
			continue
		}

		progCback("Extracting %s...\n", dstPath)

		dstFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY, f.FileInfo().Mode())
		if err != nil {
			return errors.Wrap(err, "Cannot create file")
		}
		_, err = io.Copy(dstFile, tarFile)

		if err != nil {
			return errors.Wrap(err, "Cannot write data")
		}

	}
	return nil
}
