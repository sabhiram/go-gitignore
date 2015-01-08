/*
ignore is a library which returns a new ignorer object which can
test against various paths. This is particularly useful when trying
to filter files based on a .gitignore document

The rules for parsing the input file are the same as the ones listed
in the Git docs here: http://git-scm.com/docs/gitignore

The summarized version of the same has been copied here:

    1. A blank line matches no files, so it can serve as a separator
       for readability.
    2. A line starting with # serves as a comment. Put a backslash ("\")
       in front of the first hash for patterns that begin with a hash.
    3. Trailing spaces are ignored unless they are quoted with backslash ("\").
    4. An optional prefix "!" which negates the pattern; any matching file
       excluded by a previous pattern will become included again. It is not
       possible to re-include a file if a parent directory of that file is
       excluded. Git doesn’t list excluded directories for performance reasons,
       so any patterns on contained files have no effect, no matter where they
       are defined. Put a backslash ("\") in front of the first "!" for
       patterns that begin with a literal "!", for example, "\!important!.txt".
    5. If the pattern ends with a slash, it is removed for the purpose of the
       following description, but it would only find a match with a directory.
       In other words, foo/ will match a directory foo and paths underneath it,
       but will not match a regular file or a symbolic link foo (this is
       consistent with the way how pathspec works in general in Git).
    6. If the pattern does not contain a slash /, Git treats it as a shell glob
       pattern and checks for a match against the pathname relative to the
       location of the .gitignore file (relative to the toplevel of the work
       tree if not from a .gitignore file).
    7. Otherwise, Git treats the pattern as a shell glob suitable for
       consumption by fnmatch(3) with the FNM_PATHNAME flag: wildcards in the
       pattern will not match a / in the pathname. For example,
       "Documentation/*.html" matches "Documentation/git.html" but not
       "Documentation/ppc/ppc.html" or "tools/perf/Documentation/perf.html".
    8. A leading slash matches the beginning of the pathname. For example,
       "/*.c" matches "cat-file.c" but not "mozilla-sha1/sha1.c".
    9. Two consecutive asterisks ("**") in patterns matched against full
       pathname may have special meaning:
        i.   A leading "**" followed by a slash means match in all directories.
             For example, "** /foo" matches file or directory "foo" anywhere,
             the same as pattern "foo". "** /foo/bar" matches file or directory
             "bar" anywhere that is directly under directory "foo".
        ii.  A trailing "/**" matches everything inside. For example, "abc/**"
             matches all files inside directory "abc", relative to the location
             of the .gitignore file, with infinite depth.
        iii. A slash followed by two consecutive asterisks then a slash matches
             zero or more directories. For example, "a/** /b" matches "a/b",
             "a/x/b", "a/x/y/b" and so on.
        iv.  Other consecutive asterisks are considered invalid. */
package ignore

import (
    "fmt"
    "strings"
    "regexp"
    "io/ioutil"
)

// An IgnoreParser is an interface which exposes two methods:
//   IncludesPath() - Returns true if the path will be included
//   IgnoresPath()  - Returns true if the path will be excluded
type IgnoreParser interface {
    IncludesPath(f string) bool
    IgnoresPath(f string) bool
}

// GitIgnore is a struct which contains a slice of regexp.Regexp
// patterns.
type GitIgnore struct {
    patterns []*regexp.Regexp
}

// This function pretty much attempts to mimic the parsing
// rules listed above at the start of this file
func getPatternFromLine(line string) *regexp.Regexp {
    // Strip comments
    r := regexp.MustCompile("^(.*?)#.*$")
    line = r.ReplaceAllString(line, "$1")

    // Trim string
    line = strings.Trim(line, " ")

    // Exit for no-ops
    if line == "" { return nil }

    // Temporary regex
    expr := "^" + line + "(|/.*)$"
    fmt.Printf("Line: %s has pattern: %s\n", line, expr)
    pattern, error := regexp.Compile(expr)
    if error == nil {
        return pattern
    }
    return nil
}

// Accepts a variadic set of strings, and returns a GitIgnore object which
// converts and appends the lines in the input to regexp.Regexp patterns
// held within the GitIgnore objects "patterns" field
func CompileIgnoreLines(lines ...string) (*GitIgnore, error) {
    g := new(GitIgnore)
    for _, line := range lines {
        pattern := getPatternFromLine(line)
        if pattern != nil {
            g.patterns = append(g.patterns, pattern)
        }
    }
    return g, nil
}

// Accepts a ignore file as the input, parses the lines out of the file
// and invokes the CompileIgnoreLines method
func CompileIgnoreFile(fpath string) (*GitIgnore, error) {
    buffer, error := ioutil.ReadFile(fpath)
    if error == nil {
        s := strings.Split(string(buffer), "\n")
        return CompileIgnoreLines(s...)
    }
    return nil, error
}

// IncludesPath is an interface function for the IgnoreParser interface.
// It returns true if the given GitIgnore structure would not reject the path
// being queried against
func (g GitIgnore) IncludesPath(f string) bool {
    for _, pattern := range g.patterns {
        if pattern.MatchString(f) { return false }
    }
    return true
}

// IgnoresPath is an interface function for the IgnoreParser interface.
// It returns true if the given GitIgnore structure would reject the path
// being queried against
func (g GitIgnore) IgnoresPath(f string) bool {
    return !g.IncludesPath(f)
}
