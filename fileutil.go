package goup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/karrick/godirwalk"
)

const (
	BufSize = 10 * 1024
)

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return errors.Wrap(err, "Cannot open src")
	}
	defer srcFile.Close()

	srcAttr, err := srcFile.Stat()
	if err != nil {
		return errors.Wrap(err, "Cannot get src attributes")
	}
	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("File %s alread exists", dst)
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, srcAttr.Mode())
	if err != nil {
		err = os.MkdirAll(filepath.Dir(dst), 0666)
		if err != nil {
			return errors.Wrap(err, "Error when creating parent directory in target location")
		}
		dstFile, err = os.Create(dst)
		if err != nil {
			return err
		}
	}
	defer dstFile.Close()

	buf := make([]byte, BufSize)
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := dstFile.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}

func RecursiveCopyDir(src, dst string) error {
	buf := make([]byte, BufSize)
	baseDirLen := len(src)
	err := godirwalk.Walk(src, &godirwalk.Options{
		ScratchBuffer: buf,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsDir() {
				return nil
			}
			newPath := filepath.Join(dst, osPathname[baseDirLen:])
			err := copyFile(osPathname, newPath)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			return err
		},
		PostChildrenCallback: func(osPathname string, de *godirwalk.Dirent) error {
			deChildren, err := godirwalk.ReadDirents(osPathname, buf)
			if err != nil {
				return err
			}

			if len(deChildren) > 0 {
				return nil
			}

			return os.MkdirAll(filepath.Join(dst, osPathname[baseDirLen:]), 0666)
		},
	})

	return err
}
