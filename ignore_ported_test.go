package ignore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompileIgnoreFile(t *testing.T) {
	t.Parallel()
	// TODO Add benchmarking.
	type Test struct {
		name         string
		lines        []string
		wantMatch    []string
		wantNotMatch []string
	}
	for _, tt := range []Test{
		// TODO Consider a test fuzzer or something similar, it could go a long way to getting this a bit better.
		{
			name:         "Should match simple, from root and anywhere",
			lines:        []string{"foo", "**/foo", "/**/foo"},
			wantMatch:    []string{"foo", "foo/", "/foo"},
			wantNotMatch: []string{"fooo", "ofoo"},
		},
		{
			name:         "Should match simple directories",
			lines:        []string{"foo/"},
			wantMatch:    []string{"foo/a", "/foo/", "/foo/a"},
			wantNotMatch: []string{"foo", "/foo"},
		},
		{
			name:         "Should match root extension",
			lines:        []string{"/.js", "/*.js"},
			wantMatch:    []string{".js", ".js/", ".js/a", "/hello.js"},
			wantNotMatch: []string{".jsa", ".go", ".ts"},
		},
		{
			name:         "Should ONLY match root extension ",
			lines:        []string{"/.js"},
			wantMatch:    []string{".js", ".js/", ".js/a"},
			wantNotMatch: []string{".jsa", ".go", ".ts", "main.js"},
		},
		{
			name:         "Should match extension",
			lines:        []string{".js", "*.js"},
			wantMatch:    []string{".js", ".js/", ".js/a", "a.js", "main.js"},
			wantNotMatch: []string{".jsa", "main.ts", "a.ts", "data.json", "src/main.go"},
		},
		{
			name:         "Should match wildcard extensions",
			lines:        []string{".js*", "*.j*"},
			wantMatch:    []string{".json", ".jsonnet/", ".json/a", "a.js", "main.js", "main.jsonnet"},
			wantNotMatch: []string{".ajs", "main.ts", "a.ts", "src/main.go"},
		},
		{
			name:         "Should match double wildcard directories",
			lines:        []string{"foo/**/", "foo/baz/**/"},
			wantMatch:    []string{"foo/", "foo/abc/", "foo/baz/a/", "foo/main/baz/"},
			wantNotMatch: []string{".jsa", "main.ts", "src/main.go", "foo", "/foo", "/foo/baz/a"},
		},
		{
			name:         "Should match double wildcard files",
			lines:        []string{"foo/**/*.baz", "foo/**/*.bar"},
			wantMatch:    []string{"foo/hello.baz", "foo/abc/bad.bar", "foo/baz/a/b.baz", "foo/main/baz/foo.bar"},
			wantNotMatch: []string{".jsa", "main.ts", "src/main.go", "foo", "/foo", "/foo/baz/a"},
		},
		{
			name:         "Should not match comments",
			lines:        []string{`#hello`},
			wantNotMatch: []string{"#hello"},
		},
		{
			name:      "Should match escaped comments",
			lines:     []string{`\#hello`},
			wantMatch: []string{"#hello"},
		},
		{
			name:         "Should filter exception paths",
			lines:        []string{"abc", "!abc/b"},
			wantMatch:    []string{"abc/a.js"},
			wantNotMatch: []string{"abc/b/b.js"},
		},
		{
			name:         "Should filter exception paths",
			lines:        []string{"abc", "!abc/b", "#e", `\#f`},
			wantMatch:    []string{"abc/a.js", "#f"},
			wantNotMatch: []string{"abc/b/b.js", "#e"},
		},
		{
			name:      "Should escape regex metecharacters",
			lines:     []string{"*.js", `!\*.js`, "!a#b.js", "!?.js", "#abc", `\#abc`},
			wantMatch: []string{"abc.js", "#abc"},
			wantNotMatch: []string{
				"?.js", "*.js", "a#b.js", "abc",
			},
		},
		{
			name:         "Should match wildcard directories",
			lines:        []string{"abc/*", "foo/*", "baz/bar/*"},
			wantMatch:    []string{"abc/foo", "foo/hello", "baz/bar/bat"},
			wantNotMatch: []string{"bat/bar", "baz/abc", "oof/foo"},
		},
		{
			name:         "Should match .DS_Store and other hidden files",
			lines:        []string{".DS_Store", ".d"},
			wantMatch:    []string{".DS_Store", "abc/.DS_Store", "abc/.config/.DS_Store", "/root/to/dir/.d"},
			wantNotMatch: []string{"meme.DS_Store", "meme.d", "abc/meme.d"},
		},
		{
			name:         "",
			lines:        []string{},
			wantMatch:    []string{},
			wantNotMatch: []string{},
		},
		{
			name:         "",
			lines:        []string{},
			wantMatch:    []string{},
			wantNotMatch: []string{},
		},
		{
			name:         "",
			lines:        []string{},
			wantMatch:    []string{},
			wantNotMatch: []string{},
		},
		{
			name:      "Should match pattern once",
			lines:     []string{"node_modules/"},
			wantMatch: []string{"node_modules/gulp/node_modules/abc.md", "node_modules/gulp/node_modules/abc.json"},
		},
		{
			name:      "Should match pattern twice",
			lines:     []string{"node_modules/", "node_modules/"},
			wantMatch: []string{"node_modules/gulp/node_modules/abc.md", "node_modules/gulp/node_modules/abc.json"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			obj := CompileIgnoreLines(tt.lines...)
			for _, w := range tt.wantMatch {
				assert.Equal(t, true, obj.MatchesPath(w), w+": should match")
			}
			for _, wN := range tt.wantNotMatch {
				assert.NotEqual(t, true, obj.MatchesPath(wN), wN+": should not match")
			}
		})
	}
}
