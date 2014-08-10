package gen

import "os"

// Overwrite and create new file
func createFile(filename string) (f *os.File) {
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
