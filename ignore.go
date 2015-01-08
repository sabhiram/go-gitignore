// `ignore` is a library which returns a new ignorer object
// which can test against various paths. This is particularly
// useful when trying to filter files based on a `.gitignore`
// document
package ignore

import (
    "fmt"
    "strings"

    "regexp"

    "io/ioutil"
)

// TODO: Use fmt
var _ = fmt.Println

type IgnoreParser interface {
    AcceptsPath(f string) bool
    IgnoresPath(f string) bool
}

type GitIgnore struct {
    patterns []*regexp.Regexp
}

func CompileLines(lines ...string) (*GitIgnore, error) {
    g := new(GitIgnore)
    for _, line := range lines {
        // fmt.Println("Line: " + line)
        // TODO: This is temporary:
        pattern, _ := regexp.Compile("^" + line + "$")
        g.patterns = append(g.patterns, pattern)
    }
    return g, nil
}

func CompileFile(fpath string) (*GitIgnore, error) {
    buffer, error := ioutil.ReadFile(fpath)
    if error == nil {
        s := strings.Split(string(buffer), "\n")
        fmt.Println(s)
        return CompileLines(s...)
    }
    return nil, error
}

func (g GitIgnore) AcceptsPath(f string) bool {
    for _, pattern := range g.patterns {
        if pattern.MatchString(f) { return false }
    }
    return true
}

func (g GitIgnore) IgnoresPath(f string) bool {
    return !g.AcceptsPath(f)
}
