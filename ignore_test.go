// Implement tests for the `ignore` library
package ignore

import (
    "os"

    "io/ioutil"
    "path/filepath"

    "fmt"
    "testing"

    "github.com/stretchr/testify/assert"
)

const (
    TEST_DIR = "test_fixtures"
)

// Helper function to setup a test fixture dir and write to
// it a file with the name "fname" and content "content"
func writeFileToTestDir(fname, content string) {
    testDirPath := "." + string(filepath.Separator) + TEST_DIR
    testFilePath := testDirPath + string(filepath.Separator) + fname

    _ = os.MkdirAll(testDirPath, 0755)
    _ = ioutil.WriteFile(testFilePath, []byte(content), os.ModePerm)
}

func cleanupTestDir() {
    _ = os.RemoveAll(fmt.Sprintf(".%s%s", string(filepath.Separator), TEST_DIR))
}

// Validate "CompileIgnoreLines()"
func TestCompileIgnoreLines(test *testing.T) {
    lines := []string{"abc/def", "a/b/c", "b"}
    object, error := CompileIgnoreLines(lines...)
    assert.Nil(test, error, "error from CompileIgnoreLines should be nil")

    // IncludesPath
    // Paths which should not be ignored
    assert.Equal(test, true, object.IncludesPath("abc"), "abc should not be ignored")
    assert.Equal(test, true, object.IncludesPath("def"), "def should not be ignored")
    assert.Equal(test, true, object.IncludesPath("bd"),  "bd should not be ignored")
    // Paths which should be ignored
    assert.Equal(test, false, object.IncludesPath("abc/def/child"), "abc/def/child should be ignored")
    assert.Equal(test, false, object.IncludesPath("a/b/c/d"),       "a/b/c/d should be ignored")

    object, error = CompileIgnoreLines("abc/def", "a/b/c", "b")
    assert.Nil(test, error, "error from CompileIgnoreLines should be nil")

    // IgnorePath
    assert.Equal(test, false, object.IgnoresPath("abc"), "abc should not be ignored")
    assert.Equal(test, false, object.IgnoresPath("def"), "def should not be ignored")
    assert.Equal(test, false, object.IgnoresPath("bd"),  "bd should not be ignored")

    // Paths which should be ignored
    assert.Equal(test, true, object.IgnoresPath("abc/def/child"), "abc/def/child should be ignored")
    assert.Equal(test, true, object.IgnoresPath("a/b/c/d"),       "a/b/c/d should be ignored")
}

// Validate the invalid files
func TestCompileIgnoreFile_InvalidFile(test *testing.T) {
    object, error := CompileIgnoreFile("./test_fixtures/invalid.file")
    assert.Nil(test, object, "object should be nil")
    assert.NotNil(test, error, "error should be unknown file / dir")
}

// Validate the an empty files
func TestCompileIgnoreLines_EmptyFile(test *testing.T) {
    writeFileToTestDir("test.gitignore", ``)
    defer cleanupTestDir()

    object, error := CompileIgnoreFile("./test_fixtures/test.gitignore")
    assert.Nil(test, error, "error should be nil")
    assert.NotNil(test, object, "object should not be nil")

    assert.Equal(test, false, object.IgnoresPath("a"),       "should accept all paths")
    assert.Equal(test, false, object.IgnoresPath("a/b"),     "should accept all paths")
    assert.Equal(test, false, object.IgnoresPath(".foobar"), "should accept all paths")
}

//
// FOLDER based path checking tests
//

// Validate the correct handling of the negation operator "!"
func TestCompileIgnoreLines_HandleIncludePattern(test *testing.T) {
    writeFileToTestDir("test.gitignore", `
# exclude everything except directory foo/bar
/*
!/foo
/foo/*
!/foo/bar
`)
    defer cleanupTestDir()

    object, error := CompileIgnoreFile("./test_fixtures/test.gitignore")
    assert.Nil(test, error, "error should be nil")
    assert.NotNil(test, object, "object should not be nil")

    assert.Equal(test, true,  object.IgnoresPath("a"),        "a should be ignored")
    assert.Equal(test, true,  object.IgnoresPath("foo/baz"), "foo/baz should be ignored")
    assert.Equal(test, false, object.IgnoresPath("foo"),      "foo should not be ignored")
    assert.Equal(test, false, object.IgnoresPath("/foo/bar"), "/foo/bar should not be ignored")
}

// Validate the correct handling of comments and empty lines
func TestCompileIgnoreLines_HandleSpaces(test *testing.T) {
    writeFileToTestDir("test.gitignore", `
#
# A comment

# Another comment


    # Invalid Comment

abc/def
`)
    defer cleanupTestDir()

    object, error := CompileIgnoreFile("./test_fixtures/test.gitignore")
    assert.Nil(test, error, "error should be nil")
    assert.NotNil(test, object, "object should not be nil")

    assert.Equal(test, 2, len(object.patterns), "should have two regex pattern")
    assert.Equal(test, false, object.IgnoresPath("abc/abc"), "/abc/abc should not be ignored")
    assert.Equal(test, true,  object.IgnoresPath("abc/def"), "/abc/def should be ignored")
}

// Validate the correct handling of leading / chars
func TestCompileIgnoreLines_HandleLeadingSlash(test *testing.T) {
    writeFileToTestDir("test.gitignore", `
/a/b/c
d/e/f
/g
`)
    defer cleanupTestDir()

    object, error := CompileIgnoreFile("./test_fixtures/test.gitignore")
    assert.Nil(test, error, "error should be nil")
    assert.NotNil(test, object, "object should not be nil")

    assert.Equal(test, 3, len(object.patterns), "should have 3 regex patterns")
    assert.Equal(test, true,  object.IgnoresPath("a/b/c"),   "a/b/c should be ignored")
    assert.Equal(test, true,  object.IgnoresPath("a/b/c/d"), "a/b/c/d should be ignored")
    assert.Equal(test, true,  object.IgnoresPath("d/e/f"),   "d/e/f should be ignored")
    assert.Equal(test, true,  object.IgnoresPath("g"),       "g should be ignored")
}

//
// FILE based path checking tests
//

// Validate the correct handling of files starting with # or !
func TestCompileIgnoreLines_HandleLeadingSpecialChars(test *testing.T) {
    writeFileToTestDir("test.gitignore", `
# Comment
\#file.txt
\!file.txt
file.txt
`)
    defer cleanupTestDir()

    object, error := CompileIgnoreFile("./test_fixtures/test.gitignore")
    assert.Nil(test, error, "error should be nil")
    assert.NotNil(test, object, "object should not be nil")

    assert.Equal(test, true,  object.IgnoresPath("#file.txt"),   "#file.txt should be ignored")
    assert.Equal(test, true,  object.IgnoresPath("!file.txt"),   "!file.txt should be ignored")
    assert.Equal(test, true,  object.IgnoresPath("a/!file.txt"), "a/!file.txt should be ignored")
    assert.Equal(test, true,  object.IgnoresPath("file.txt"),    "file.txt should be ignored")
    assert.Equal(test, true,  object.IgnoresPath("a/file.txt"),  "a/file.txt should be ignored")
    assert.Equal(test, false, object.IgnoresPath("file2.txt"),   "file2.txt should not be ignored")

}

// Validate the correct handling matching files only within a given folder
func TestCompileIgnoreLines_HandleAllFilesInDir(test *testing.T) {
    writeFileToTestDir("test.gitignore", `
Documentation/*.html
`)
    defer cleanupTestDir()

    object, error := CompileIgnoreFile("./test_fixtures/test.gitignore")
    assert.Nil(test, error, "error should be nil")
    assert.NotNil(test, object, "object should not be nil")

    assert.Equal(test, true,  object.IgnoresPath("Documentation/git.html"),             "Documentation/git.html should be ignored")
    assert.Equal(test, false, object.IgnoresPath("Documentation/ppc/ppc.html"),         "Documentation/ppc/ppc.html should not be ignored")
    assert.Equal(test, false, object.IgnoresPath("tools/perf/Documentation/perf.html"), "tools/perf/Documentation/perf.html should not be ignored")
}

// Validate the correct handling of "**"
func TestCompileIgnoreLines_HandleDoubleStar(test *testing.T) {
    writeFileToTestDir("test.gitignore", `
**/foo
bar
`)
    defer cleanupTestDir()

    object, error := CompileIgnoreFile("./test_fixtures/test.gitignore")
    assert.Nil(test, error, "error should be nil")
    assert.NotNil(test, object, "object should not be nil")

    assert.Equal(test, true,  object.IgnoresPath("foo"),     "foo should be ignored")
    assert.Equal(test, true,  object.IgnoresPath("baz/foo"), "baz/foo should be ignored")
    assert.Equal(test, true,  object.IgnoresPath("bar"),     "bar should be ignored")
    assert.Equal(test, true,  object.IgnoresPath("baz/bar"), "baz/bar should be ignored")
}

// Validate the correct handling of leading slash
func TestCompileIgnoreLines_HandleLeadingSlashPath(test *testing.T) {
    writeFileToTestDir("test.gitignore", `
/*.c
`)
    defer cleanupTestDir()

    object, error := CompileIgnoreFile("./test_fixtures/test.gitignore")
    assert.Nil(test, error, "error should be nil")
    assert.NotNil(test, object, "object should not be nil")

    assert.Equal(test, true,  object.IgnoresPath("hello.c"),     "hello.c should be ignored")
    assert.Equal(test, false, object.IgnoresPath("foo/hello.c"), "foo/hello.c should not be ignored")
}

func ExampleCompileIgnoreLines() {
    ignoreObject, error := CompileIgnoreLines([]string{"node_modules", "*.out", "foo/*.c"}...)
    if error != nil {
        panic("Error when compiling ignore lines: " + error.Error())
    }

    fmt.Println(ignoreObject.IgnoresPath("node_modules/test/foo.js"))
    fmt.Println(ignoreObject.IgnoresPath("node_modules2/test.out"))
    fmt.Println(ignoreObject.IgnoresPath("test/foo.js"))

    // Output:
    // true
    // true
    // false
}
