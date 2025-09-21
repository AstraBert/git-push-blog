package commons

import "testing"

func TestPathExists(t *testing.T) {
	testCases := []struct {
		file   string
		exists bool
	}{
		{"../test_files/1.md", true},
		{"../test_files/2.md", true},
		{"../test_files/3.md", false},
		{"../test_files/file.txt", false},
	}
	for _, tc := range testCases {
		if tc.exists != PathExists(tc.file) {
			t.Errorf("Expecting PathExists to return %v, but returned %v instead", tc.exists, PathExists(tc.file))
		}
	}
}
