package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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
	file, err := processChangelogFile(opt)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	err = saveChangelogFile(opt.Filename, file)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	if opt.RubyLib {
		found, filename, err := existsRubyLibVersionFile()
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
		if !found {
			fmt.Println("No ruby lib version file found")
			os.Exit(0)
		}

		buf, permissions, err := processRubyLibVersionFile(filename, file.Version)
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
		err = os.WriteFile(filename, buf, permissions)
		if err != nil {
			fmt.Println(err)
			os.Exit(5)
		}
	}

}

func processChangelogFile(opts Options) (ChangelogFile, error) {
	lines, err := os.ReadFile(opts.Filename)
	if err != nil {
		return ChangelogFile{}, err
	}
	file, err := parseChangelog(lines)
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

func (s *SemverLine) VersionWithPrefix() string {
	return fmt.Sprintf("%s%s", s.Prefix, s.String())
}

func (s *SemverLine) String() string {
	return fmt.Sprintf("%d.%d.%d", s.Major, s.Minor, s.Patch)
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

type ChangelogFile struct {
	Header  []string
	Version SemverLine
	Comment string
	Body    []string
}

func (file ChangelogFile) Lines() []string {
	lines := []string{}
	lines = append(lines, file.Header...)
	lines = append(lines, file.Version.VersionWithPrefix())
	lines = append(lines, file.Comment)
	lines = append(lines, file.Body...)
	return lines
}

func parseChangelog(content []byte) (ChangelogFile, error) {
	lines := strings.Split(string(content), "\n")
	file := ChangelogFile{Header: []string{}, Body: []string{}}

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
	}
	return file, nil
}

func saveChangelogFile(filename string, file ChangelogFile) error {
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

var (
	semverRegex        = regexp.MustCompile(`(?P<Prefix>.*)(?P<Major>\d+)\.(?P<Minor>\d+)\.(?P<Patch>\d+)`)
	rubyLibSemverRegex = regexp.MustCompile(`(?P<Prefix>\s*VERSION\s*=\s*['"]{1})(?P<Major>\d+)\.(?P<Minor>\d+)\.(?P<Patch>\d+)(?P<Suffix>['"]{1})(?P<Freeze>\.freeze)?`)
)

func parseSemver(line string) *SemverLine {

	match := semverRegex.FindStringSubmatch(line)
	verMap := make(map[string]string)
	for i, name := range semverRegex.SubexpNames() {
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

func processRubyLibVersionFile(filename string, version SemverLine) ([]byte, fs.FileMode, error) {
	permissions := fs.FileMode(0)
	buf := []byte{}

	// Read the permissions so we can write them back
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return buf, permissions, err
	}
	permissions = fileInfo.Mode().Perm()

	file, err := os.ReadFile(filename)
	if err != nil {
		return buf, permissions, err
	}
	lines := strings.Split(string(file), "\n")
	for i, line := range lines {
		// Look for a matching line
		match := rubyLibSemverRegex.FindStringSubmatch(line)
		verMap := make(map[string]string)
		for i, name := range rubyLibSemverRegex.SubexpNames() {
			if i > 0 && i <= len(match) {
				verMap[name] = match[i]
			}
		}
		if len(verMap) == 0 {
			continue
		}
		lines[i] = fmt.Sprintf("%s%s%s", verMap["Prefix"], version.String(), verMap["Suffix"])
		if verMap["Freeze"] != "" {
			lines[i] += verMap["Freeze"]
		}
		break
	}
	return []byte(strings.Join(lines, "\n")), permissions, nil
}

func existsRubyLibVersionFile() (bool, string, error) {
	versionFiles, err := filepath.Glob("lib/*/version.rb")
	if err != nil {
		return false, "", err
	}

	if len(versionFiles) == 0 {
		return false, "", nil
	} else if len(versionFiles) == 1 {
		return true, versionFiles[0], nil
	}
	return false, "", fmt.Errorf("Found multiple 'lib/*/version.rb' files: %s", strings.Join(versionFiles, ", "))
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
	RubyLib  bool
}

func parseOpts() (Options, error) {
	var opts Options
	flag.StringVar(&opts.Filename, "filename", "", "CHANGELOG.md filename to read")
	flag.StringVar(&opts.Comment, "comment", "", "comment to add into the CHANGELOG.md")
	flag.BoolVar(&opts.Major, "major", false, "Specify this flag to bump the major version")
	flag.BoolVar(&opts.Minor, "minor", false, "Specify this flag to bump the minor version")
	flag.BoolVar(&opts.Patch, "patch", false, "Specify this flag to bump the patch version")
	flag.BoolVar(&opts.RubyLib, "ruby-lib", true, "Specify this flag to automatically update the version in a ruby lib version file 'lib/*/version.rb'")
	flag.Parse()
	if opts.Filename == "" {
		return opts, fmt.Errorf("filename is required")
	}
	if opts.Comment == "" {
		return opts, fmt.Errorf("comment is required")
	}
	// Convert esacpes to newlines in the comment
	opts.Comment = strings.ReplaceAll(opts.Comment, `\n`, "\n")

	if (opts.Major && (opts.Minor || opts.Patch)) ||
		(opts.Minor && (opts.Major || opts.Patch)) ||
		(opts.Patch && (opts.Major || opts.Minor)) ||
		(!opts.Minor && !opts.Major && !opts.Patch) {
		return opts, fmt.Errorf("exactly one of major, minor and patch must be selected")
	}
	return opts, nil
}
