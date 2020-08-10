package fileutils

import (
	"testing"
)

type fileOrDirExpect struct {
	path        string
	expected    bool
	errExpected bool
}

func TestIsFile(t *testing.T) {
	expectedResults := []fileOrDirExpect{
		{"/", false, false},
		{"/etc", false, false},
		{"/etc/passwd", true, false},
		{"/asdf", false, true},
	}

	for _, test := range expectedResults {
		isFile, err := IsFile(test.path)
		t.Logf("Path %s is a file?  %v; err is '%v'; err expected? %v\n",
			test.path, isFile, err, test.errExpected)

		if test.errExpected && err == nil ||
			!test.errExpected && err != nil ||
			isFile != test.expected {
			t.Fail()
		}
	}
}
func TestIsDir(t *testing.T) {
	expectedResults := []fileOrDirExpect{
		{"/", true, false},
		{"/etc", true, false},
		{"/etc/passwd", false, false},
		{"/asdf", false, true},
		{"/asdf/", false, true},
	}

	for _, test := range expectedResults {
		isDir, err := IsDir(test.path)
		t.Logf("Path %s is a dir?  %v; err is '%v'; err expected? %v\n",
			test.path, isDir, err, test.errExpected)

		if test.errExpected && err == nil ||
			!test.errExpected && err != nil ||
			isDir != test.expected {
			t.Logf("\tERROR")
			t.Fail()
		}
	}
}

func TestPrettySize(t *testing.T) {
	a := 21.03 * 1024 * 1024
	a1 := int64(a)
	b := 520.42 * 1024 * 1024 * 1024 * 1024
	b1 := int64(b)

	vals := []int64{
		10,
		15 * 1024,
		a1,
		520 * 1024 * 1024 * 1024,
		b1,
	}

	expected := []string{
		"10B",
		"15.00KB",
		"21.03MB",
		"520.00GB",
		"520.42TB",
	}

	for i, val := range vals {
		exp := expected[i]
		actual := PrettySize(val)
		t.Logf("input=%d\toutput=%s; expected=%s\n", val, actual, exp)
		if actual != exp {
			t.Fail()
		}
	}
}

func TestStripOffMBSuffix(t *testing.T) {
	expected := map[string]string{
		"50":    "50",
		"60mb":  "60",
		"72Mb":  "72",
		"500MB": "500",
	}

	for input, expected := range expected {
		actual := stripOffMBSuffix(input)
		t.Logf("Input: %s; Returned: %s; Expected: %s\n",
			input, actual, expected)
		if actual != expected {
			t.Fail()
		}
	}
}
