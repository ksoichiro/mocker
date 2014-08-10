package gen

import (
	"os"
	"path/filepath"
)

// Overwrite and create new file
func createFile(filename string) (f *os.File) {
	dir := filepath.Dir(filename)
	if !fileExists(dir) {
		os.MkdirAll(dir, 0777)
	}
	if fileExists(filename) {
		os.Remove(filename)
	}
	f, _ = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	return
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
