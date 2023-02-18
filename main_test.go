package main

import (
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
