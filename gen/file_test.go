package gen

import "testing"

func TestFileExists(t *testing.T) {
	var testcases = []struct {
		filename string
		expect   bool
	}{
		{"../main.go", true},
		{"../gen", true},
		{"../foo-bar", false},
	}
	for _, tc := range testcases {
		if actual := fileExists(tc.filename); actual != tc.expect {
			t.Errorf("Expected %t but %t: %s", tc.expect, actual, tc.filename)
		}
	}
}
