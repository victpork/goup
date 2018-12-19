package goup

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	Format            = "zip"
	DefaultInstallDir = "C:\\Go\\bin"
)

func ExtractArchive(srcFile *os.File, size int64, targetPath string, progCback func(format string, arg ...interface{})) error {
	zipFile, err := zip.NewReader(srcFile, size)
	if err != nil {
		return err
	}

	for _, f := range zipFile.File {
		if f.FileInfo().IsDir() {
			continue
		}
		dstPath := filepath.Join(targetPath, f.FileHeader.Name[3:])
		progCback("Extracting %s...\n", dstPath)
		err = os.MkdirAll(filepath.Dir(dstPath), 0666)
		if err != nil {
			return errors.Wrap(err, "Cannot create directory")
		}
		dstFile, err := os.Create(dstPath)
		if err != nil {
			return errors.Wrap(err, "Cannot create file")
		}
		afr, err := f.Open()
		if err != nil {
			return errors.Wrap(err, "Cannot open file")
		}
		_, err = io.Copy(dstFile, afr)
		if err != nil {
			afr.Close()
			return err
		}
		afr.Close()
	}
	return nil
}
