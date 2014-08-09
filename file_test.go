package main

import "testing"

func TestFileExists(t *testing.T) {
	if !fileExists("mocker.go") {
		t.Errorf("Expected true but false")
	}
	if fileExists("foo-bar") {
		t.Errorf("Expected false but true")
	}
}
