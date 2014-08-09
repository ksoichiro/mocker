package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyTemplate(t *testing.T) {
	tmpDir := "_test_templates"
	toDir := "_test_out"
	os.MkdirAll(tmpDir, 0777)
	os.MkdirAll(toDir, 0777)

	file := "test.txt"
	filename := filepath.Join(tmpDir, file)
	f, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()
	f.WriteString("This is a test.\n")
	f.Close()

	err := copyTemplate(filename, toDir)
	if err != nil {
		fmt.Printf("%v", err)
	}

	if _, err = os.Stat(filepath.Join(toDir, file)); err != nil && os.IsNotExist(err) {
		t.Errorf("Expected copy success but fail\n")
	}

	os.RemoveAll(tmpDir)
	os.RemoveAll(toDir)
}
