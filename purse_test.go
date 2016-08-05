package purse

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	dirname  string
	fixtures = map[string]string{
		"insert.sql":        "INSERT INTO post (slug, title, created, markdown, html)\nVALUES (?, ?, ?, ?, ?)",
		"query_all.sql":     "SELECT\nid,\nslug,\ntitle,\ncreated,\nmarkdown,\nhtml\nFROM post",
		"query_by_slug.sql": "SELECT\nid,\nslug,\ntitle,\ncreated,\nmarkdown,\nhtml\nFROM post\nWHERE slug = ?",
		"foo/bar.sql":       "SELECT foo FROM bar WHERE foobar=1",
	}
)

func init() {
	dirname = filepath.Join(".", "fixtures")

	// replace newlines for running unit tests on windows
	if runtime.GOOS == "windows" {
		for k, v := range fixtures {
			fixtures[k] = strings.Replace(v, "\n", "\r\n", -1)
		}
	}
}

func TestNew(t *testing.T) {
	s, err := New(dirname)
	if err != nil {
		t.Errorf("unexpected error from New() on fixtures directory")
	}

	if len(fixtures) != len(s.files) {
		t.Errorf("invalid number of loaded SQL files")
	}

	for key, _ := range fixtures {
		_, ok := s.files[key]
		if !ok {
			t.Errorf("unable to find loaded file %s in file map", key)
		}
	}

	// verify only SQL files were loaded
	for key, _ := range s.files {
		if filepath.Ext(key) != ext {
			t.Errorf("loaded unexpected file type: %s", key)
		}
	}

	// try to load file instead of directory
	_, err = New(filepath.Join(".", "purse.go"))
	if err == nil {
		t.Errorf("expected error trying to load from non-directory")
	}

	// try to load directory that does not exist
	_, err = New(filepath.Join(".", "foo"))
	if err == nil {
		t.Errorf("expected error trying to load directory that does not exist")
	}
}

func TestGet(t *testing.T) {
	s, err := New(dirname)
	if err != nil {
		t.Errorf("unexpected error from New() on fixtures directory")
	}

	for key, val := range fixtures {
		v, ok := s.Get(key)
		if !ok {
			t.Errorf("unable to find loaded file %s in file map", key)
		}
		if v != val {
			t.Errorf("invalid %s file content:\n%v\n%v", key, []byte(v), []byte(val))
		}
	}
}

func BenchmarkGet(b *testing.B) {
	s, _ := New(dirname)
	var key string = "query_by_slug.sql"

	for i := 0; i < b.N; i++ {
		s.Get(key)
	}
}
