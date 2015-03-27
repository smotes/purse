// purse is a Go package and tool for loading and embedding SQL file contents into Go programs.
//
// To install the package, use go install:
//
// 		$ go get github.com/smotes/purse
//
// purse is also a command line tool to automate the creation of a Purse implementation given a specified
// directory of SQL files. Given the directory of SQL files and the directory of the output file,
// purse will create a Go source file (named out.go by default) which contains an implementation and
// instantiation (bound to variable gen by default) of a Purse interface driven by a map[string]string
// literal representing the SQL directory's files' contents.
//
// The purse package is meant to be used during development where SQL files are changing often
// and need to be reloaded into memory on each program execution. Contrastingly, the purse tool
// is meant to be used in production environments where the SQL files' contents can be embedded
// directly into the compiled binary, easing deployment.
//
// To install the tool, use go install:
//
//		$ go get github.com/smotes/purse/cmd/purse
// 		$ go install github.com/smotes/purse/cmd/purse
//		$ purse [args]
//
// The command syntax is:
//
//		$ purse -in="input/dir" -out="output/dir" [-file="out.go"] [-name="gen"] [-pack="main"]
//
// The input directory and output directory paths must either be absolute or relative to the
// package using it via go generate, or relative to the current working directory where the
// the command was executed.
//
// To override the default output source file name (out.go), provide the optional -file flag.
//
// To override the default variable name (gen) of the generated Purse, provide the optional -name flag.
//
// To set or override the `$GOPACKAGE` environment variable, provide the optional `-pack` flag.
//
// This process should generally be handled using go generate. Add a comment in one of your go source files,
// like so:
//
// 		//go:generate purse -in="./fixtures" -out="."
//
// Then run go generate:
//
//  	$ go generate
//
// Note that the `-pack` flag is not necessary when using go generate, as it sets the environment variable
// automatically. Refer to the [documentation](https://golang.org/cmd/go/) on the `go` command for more information.
//
package purse // import "github.com/smotes/purse"

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

const (
	ext string = ".sql"
)

// Purse is the interface representing a collection of SQL files whose
// contents are accessed via the basic Get method.
//
// Implementations of this interface should be safe for concurrent use by
// multiple goroutines.
type Purse interface {
	Get(string) (string, bool)
}

// MemoryPurse is an implementation of Purse that uses a map of strings to
// represent the contents of SQL files found within a specified directory.
// It is safe for concurrent use by multiple goroutines.
type MemoryPurse struct {
	mu    sync.RWMutex
	files map[string]string
}

// New loads SQL files' contents in the specified directory dir into memory.
// A file's contents can be accessed with the Get method.
//
// New returns an error if the directory does not exist or is not a directory.
func New(dir string) (*MemoryPurse, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fis, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	p := &MemoryPurse{files: make(map[string]string, 0)}

	for _, fi := range fis {
		if !fi.IsDir() && filepath.Ext(fi.Name()) == ext {
			f, err := os.Open(filepath.Join(dir, fi.Name()))
			if err != nil {
				return nil, err
			}
			b, err := ioutil.ReadAll(f)
			if err != nil {
				return nil, err
			}
			p.files[fi.Name()] = string(b)
		}
	}
	return p, nil
}

// Get returns a SQL file's contents as a string.
// If the file does not exist or does exists but had a read error,
// then v == "" and ok == false.
func (p *MemoryPurse) Get(filename string) (v string, ok bool) {
	p.mu.RLock()
	v, ok = p.files[filename]
	p.mu.RUnlock()
	return
}

// Files returns a slice of filenames for all loaded SQL files.
func (p *MemoryPurse) Files() []string {
	fs := make([]string, len(p.files))
	i := 0
	for k, _ := range p.files {
		fs[i] = k
		i++
	}
	return fs
}
