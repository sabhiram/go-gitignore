// Implement tests for the `ignore` library
package ignore

import (
    "fmt"

    "testing"
    "github.com/stretchr/testify/assert"
)

// Validate `CompileIgnoreLines()`
func TestCompileIgnoreLines(test *testing.T) {
    fmt.Println(" * Testing CompileIgnoreLines()")

    lines := []string{"abc/def", "a/b/c", "b"}
    object2, error := CompileIgnoreLines(lines...)

    // Validate no error
    assert.Nil(test, error, "error from CompileIgnoreLines should be nil")

    // IncludesPath
    // Paths which should not be ignored
    assert.Equal(test, object2.IncludesPath("abc"),           true,  "abc should be accepted")
    assert.Equal(test, object2.IncludesPath("def"),           true,  "def should be accepted")
    assert.Equal(test, object2.IncludesPath("bd"),            true,  "bd should be accepted")

    // Paths which should be ignored
    assert.Equal(test, object2.IncludesPath("abc/def/child"), false, "abc/def/child should be rejected")
    assert.Equal(test, object2.IncludesPath("a/b/c/d"),       false, "a/b/c/d should be rejected")

    // IgnorePath
    assert.Equal(test, object2.IgnoresPath("abc"),           false, "abc should be accepted")
    assert.Equal(test, object2.IgnoresPath("def"),           false, "def should be accepted")
    assert.Equal(test, object2.IgnoresPath("bd"),            false, "bd should be accepted")

    // Paths which should be ignored
    assert.Equal(test, object2.IgnoresPath("abc/def/child"), true,  "abc/def/child should be rejected")
    assert.Equal(test, object2.IgnoresPath("a/b/c/d"),       true,  "a/b/c/d should be rejected")
}

// Validate `CompileIgnoreFile()` for invalid files
func TestCompileIgnoreFile_InvalidFile(test *testing.T) {
    fmt.Println(" * Testing CompileIgnoreFile() for invalid file")

    object, error := CompileIgnoreFile("./test_fixtures/invalid.file")
    assert.Nil(test, object, "object should be nil")
    assert.NotNil(test, error, "error should be unknown file / dir")
}

// Validate `CompileIgnoreFile()` for an empty files
func TestCompileIgnoreLines_EmptyFile(test *testing.T) {
    fmt.Println(" * Testing CompileIgnoreFile() for empty file")

    object, error := CompileIgnoreFile("./test_fixtures/test.gitignore")
    assert.Nil(test, error, "error should be nil")
    assert.NotNil(test, object, "object should not be nil")

    assert.Equal(test, object.IgnoresPath("a"),       false, "should accept all paths")
    assert.Equal(test, object.IgnoresPath("a/b"),     false, "should accept all paths")
    assert.Equal(test, object.IgnoresPath(".foobar"), false, "should accept all paths")
}
