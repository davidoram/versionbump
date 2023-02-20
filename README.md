# versionbump
A command line tool to add a new semantic version into a CHANGELOG.md file

Download binaries from the releases

To see usage run:

```bash
$ versionbump-darwin-amd64 -help

Usage of versionbump-darwin-amd64:
  -comment string
    	comment to add into the CHANGELOG.md
  -filename string
    	CHANGELOG.md filename to read
  -major
    	Specify this flag to bump the major version
  -minor
    	Specify this flag to bump the minor version
  -patch
    	Specify this flag to bump the patch version
```

Example usage:

`versionbump-darwin-amd64 -comment "\nThis is a test\n" -filename CHANGELOG.md -minor`

Note that inside the `-comment` field `\n` is replaced with a linebreak.
