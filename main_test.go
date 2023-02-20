package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSemverString(t *testing.T) {
	var testString = []struct {
		semver SemverLine
		str    string
	}{
		{semver: SemverLine{Prefix: "##n = ", Major: 1, Minor: 2, Patch: 3}, str: "##n = 1.2.3"},
		{semver: SemverLine{Prefix: "# ", Major: 1, Minor: 2, Patch: 3}, str: "# 1.2.3"},
		{semver: SemverLine{Prefix: "", Major: 1, Minor: 2, Patch: 3}, str: "1.2.3"},
	}
	for _, test := range testString {
		assert.Equal(t, test.str, test.semver.String())
	}
}

func TestSemverIncrementMajor(t *testing.T) {
	var testString = []struct {
		semver SemverLine
		str    string
	}{
		{semver: SemverLine{Prefix: "", Major: 1, Minor: 2, Patch: 3}, str: "2.0.0"},
		{semver: SemverLine{Prefix: "", Major: 9, Minor: 0, Patch: 0}, str: "10.0.0"},
	}
	for _, test := range testString {
		test.semver.IncrementMajor()
		assert.Equal(t, test.str, test.semver.String())
	}
}

func TestSemverIncrementMinor(t *testing.T) {
	var testString = []struct {
		semver SemverLine
		str    string
	}{
		{semver: SemverLine{Prefix: "", Major: 1, Minor: 2, Patch: 3}, str: "1.3.0"},
		{semver: SemverLine{Prefix: "", Major: 9, Minor: 99, Patch: 0}, str: "9.100.0"},
	}
	for _, test := range testString {
		test.semver.IncrementMinor()
		assert.Equal(t, test.str, test.semver.String())
	}
}

func TestSemverIncrementPatch(t *testing.T) {
	var testString = []struct {
		semver SemverLine
		str    string
	}{
		{semver: SemverLine{Prefix: "", Major: 1, Minor: 2, Patch: 3}, str: "1.2.4"},
		{semver: SemverLine{Prefix: "", Major: 9, Minor: 99, Patch: 99}, str: "9.99.100"},
	}
	for _, test := range testString {
		test.semver.IncrementPatch()
		assert.Equal(t, test.str, test.semver.String())
	}
}

func TestProcessFile(t *testing.T) {
	// Find the paths of all input files in the data directory.
	paths, err := filepath.Glob(filepath.Join("testdata", "*.md"))
	if err != nil {
		t.Fatal(err)
	}

	for _, testname := range paths {

		// Each path turns into a test: the test name is the filename without the
		// extension.
		t.Run(testname, func(t *testing.T) {
			t.Logf("Testing %s", testname)
			opts := Options{
				Filename: testname,
				Comment:  "\nThis text is added to the top of the file [ticket](http://link/to/a/ticket)\n",
				Minor:    true,
			}
			outFile, err := processFile(opts)
			assert.NoError(t, err)

			// Each input file is expected to have a "golden output" file, with the
			// same path except it has a .golden extension.
			goldenfile := testname + ".golden"
			want, err := os.ReadFile(goldenfile)
			if err != nil {
				t.Fatal("error reading golden file:", err)
			}

			output := []byte(strings.Join(outFile.Lines(), "\n"))
			assert.Equal(t, want, output)
		})
	}
}
