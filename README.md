# versionbump
A command line tool to add a new semantic version into a CHANGELOG.md file, and also set the new semantic version in a ruby library's `lib/*/version.rb` file if one exists.

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
  -ruby-lib
    	Specify this flag to automatically update the version in a ruby lib version file 'lib/*/version.rb' (default true)
```

Example usage:

`versionbump-darwin-amd64 -comment "\nThis is a test\n" -filename CHANGELOG.md -minor`

Note that inside the `-comment` field `\n` is replaced with a linebreak.

## Release

Create a new release manually. Github actions will detect this, and build the artifacts and attach them to the new release.
