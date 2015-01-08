// Implement tests for the `ignore` library
package ignore

import (
    "fmt"

    "testing"
    "github.com/stretchr/testify/assert"
)

// Validate `CompileLines()`
func TestCompileLines(test *testing.T) {
    fmt.Println("")
    fmt.Println("Testing ignore.CompileLines(s ...string)")

    lines := []string{"abc/def", "a/b/c", "b"}
    object2, error := CompileLines(lines...)

    // Validate no error
    assert.Nil(test, error, "error from CompileLines should be nil")

    // AcceptsPath
    // Paths which should not be ignored
    assert.Equal(test, object2.AcceptsPath("abc"),           true,  "abc should be accepted")
    assert.Equal(test, object2.AcceptsPath("def"),           true,  "def should be accepted")
    assert.Equal(test, object2.AcceptsPath("bd"),            true,  "bd should be accepted")

    // Paths which should be ignored
    assert.Equal(test, object2.AcceptsPath("abc/def/child"), false, "abc/def/child should be rejected")
    assert.Equal(test, object2.AcceptsPath("a/b/c/d"),       false, "a/b/c/d should be rejected")

    // IgnorePath
    assert.Equal(test, object2.IgnoresPath("abc"),           false, "abc should be accepted")
    assert.Equal(test, object2.IgnoresPath("def"),           false, "def should be accepted")
    assert.Equal(test, object2.IgnoresPath("bd"),            false, "bd should be accepted")

    // Paths which should be ignored
    assert.Equal(test, object2.IgnoresPath("abc/def/child"), true,  "abc/def/child should be rejected")
    assert.Equal(test, object2.IgnoresPath("a/b/c/d"),       true,  "a/b/c/d should be rejected")
}

// Validate `CompileFile()` for invalid files
func TestCompileFile_InvalidFile(test *testing.T) {
    fmt.Println("")
    fmt.Println("Testing CompileFile() for invalid file")

    object, error := CompileFile("./test_fixtures/invalid.file")
    assert.Nil(test, object, "object should be nil")
    assert.NotNil(test, error, "error should be unknown file / dir")
}

// Validate `CompileFile()` for an empty files
func TestCompileLines_EmptyFile(test *testing.T) {
    fmt.Println("")
    fmt.Println("Testing CompileFile() for empty file")

    object, error := CompileFile("./test_fixtures/test.gitignore")
    assert.Nil(test, error, "error should be nil")
    assert.NotNil(test, object, "object should not be nil")

    assert.Equal(test, object.IgnoresPath("a"),       false, "should accept all paths")
    assert.Equal(test, object.IgnoresPath("a/b"),     false, "should accept all paths")
    assert.Equal(test, object.IgnoresPath(".foobar"), false, "should accept all paths")
}
