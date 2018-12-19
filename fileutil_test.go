package goup

import (
	"os"
	"path/filepath"
	"testing"
)

func prepareTestDirStruct(t *testing.T) {
	t.Helper()
	err := os.MkdirAll(filepath.FromSlash("test/A"), 0666)
	if err != nil {
		t.Error("Cannot create test directory", err)
	}
	err = os.MkdirAll(filepath.FromSlash("test/B"), 0666)
	if err != nil {
		t.Error("Cannot create test directory", err)
	}
	f1, err := os.Create(filepath.FromSlash("test/TestFile1"))
	if err != nil {
		t.Error("Cannot create test file", err)
	}
	f1.WriteString("HelloWorld1")
	f1.Close()
	f2, err := os.Create(filepath.FromSlash("test/A/TestFile2"))
	if err != nil {
		t.Error("Cannot create test file", err)
	}
	f2.WriteString("HelloWorld2")
	f2.Close()
}

func cleanUp(t *testing.T) {
	t.Helper()
	os.RemoveAll("test")
	os.RemoveAll("test2")
}

func TestCopyDir(t *testing.T) {
	cleanUp(t)
	prepareTestDirStruct(t)
	err := RecursiveCopyDir("test", "test2")
	if err != nil {
		t.Error(err)
	}
	if _, err := os.Stat(filepath.FromSlash("test2/TestFile1")); err != nil {
		t.Error("test2/TestFile1 not found")
	}
	if _, err := os.Stat(filepath.FromSlash("test2/A/TestFile2")); err != nil {
		t.Error("test2/A/TestFile2 not found")
	}
	if _, err := os.Stat(filepath.FromSlash("test2/B")); err != nil {
		t.Error("test2/B not found")
	}
	cleanUp(t)
}
