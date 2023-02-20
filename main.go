package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	opt, err := parseOpts()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	file, err := processFile(opt)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	err = saveFile(opt.Filename, file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func processFile(opts Options) (File, error) {
	lines, err := os.ReadFile(opts.Filename)
	if err != nil {
		return File{}, err
	}
	file, err := parse(lines)
	if err != nil {
		return file, err
	}
	if opts.Major {
		file.Version.IncrementMajor()
	} else if opts.Minor {
		file.Version.IncrementMinor()
	} else {
		file.Version.IncrementPatch()
	}
	file.Comment = opts.Comment
	return file, nil
}

type SemverLine struct {
	Prefix string
	Major  int
	Minor  int
	Patch  int
}

func (s *SemverLine) String() string {
	return fmt.Sprintf("%s%d.%d.%d", s.Prefix, s.Major, s.Minor, s.Patch)
}

func (s *SemverLine) IncrementMajor() {
	s.Major++
	s.Minor = 0
	s.Patch = 0
}

func (s *SemverLine) IncrementMinor() {
	s.Minor++
	s.Patch = 0
}

func (s *SemverLine) IncrementPatch() {
	s.Patch++
}

type File struct {
	Header  []string
	Version SemverLine
	Comment string
	Body    []string
}

func (file File) Lines() []string {
	lines := []string{}
	lines = append(lines, file.Header...)
	lines = append(lines, file.Version.String())
	lines = append(lines, file.Comment)
	lines = append(lines, file.Body...)
	return lines
}

func parse(content []byte) (File, error) {
	lines := strings.Split(string(content), "\n")
	file := File{Header: []string{}, Body: []string{}}

	findingVersion := true
	for _, line := range lines {
		// If we are still looking for the version we are in the header
		if findingVersion {
			if semver := parseSemver(line); semver != nil {
				file.Version = *semver
				findingVersion = false

				// Add the line to the body,so it's retained it in the new version of the file
				file.Body = append(file.Body, line)
			} else {
				file.Header = append(file.Header, line)
			}
		} else {
			file.Body = append(file.Body, line)
		}
		fmt.Println(line)
	}
	return file, nil
}

func saveFile(filename string, file File) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, line := range file.Lines() {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func parseSemver(line string) *SemverLine {
	r := regexp.MustCompile(`(?P<Prefix>.*)(?P<Major>\d+)\.(?P<Minor>\d+)\.(?P<Patch>\d+)`)
	match := r.FindStringSubmatch(line)
	verMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			verMap[name] = match[i]
		}
	}
	if len(verMap) == 0 {
		return nil
	}
	return &SemverLine{
		Prefix: verMap["Prefix"],
		Major:  mustAtoi(verMap["Major"]),
		Minor:  mustAtoi(verMap["Minor"]),
		Patch:  mustAtoi(verMap["Patch"])}
}

func mustAtoi(s string) int {
	ver, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return ver
}

type Options struct {
	Filename string
	Comment  string
	Major    bool
	Minor    bool
	Patch    bool
}

func parseOpts() (Options, error) {
	var opts Options
	flag.StringVar(&opts.Filename, "filename", "", "CHANGELOG.md filename to read")
	flag.StringVar(&opts.Comment, "comment", "", "comment to add into the CHANGELOG.md")
	flag.BoolVar(&opts.Major, "major", false, "Specify this flag to bump the major version")
	flag.BoolVar(&opts.Minor, "minor", false, "Specify this flag to bump the minor version")
	flag.BoolVar(&opts.Patch, "patch", false, "Specify this flag to bump the patch version")
	flag.Parse()
	if opts.Filename == "" {
		return opts, fmt.Errorf("filename is required")
	}
	if opts.Comment == "" {
		return opts, fmt.Errorf("comment is required")
	}
	if (opts.Major && (opts.Minor || opts.Patch)) ||
		(opts.Minor && (opts.Major || opts.Patch)) ||
		(opts.Patch && (opts.Major || opts.Minor)) ||
		(!opts.Minor && !opts.Major && !opts.Patch) {
		return opts, fmt.Errorf("major, minor and patch are mutually exclusive")
	}
	return opts, nil
}
