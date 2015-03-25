package purse

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

const (
	ext string = ".sql"
)

// Purse is a key/value collection of loaded SQL files by name : content.
// It is safe for concurrent use by multiple goroutines.
type Purse struct {
	mu    sync.RWMutex
	files map[string]string
}

// New loads all SQL files in the specified directory dir.
//
// A loaded file's contents can be accessed via Get().
//
// Returns an error if the directory does not exist or is not a directory.
func New(dir string) (*Purse, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fis, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	s := &Purse{files: make(map[string]string, 0)}

	for _, fi := range fis {
		if !fi.IsDir() && filepath.Ext(fi.Name()) == ext {
			f, err := os.Open(filepath.Join(dir, fi.Name()))
			if err != nil {
				return nil, err
			}
			defer f.Close()

			fd, err := ioutil.ReadAll(f)
			if err != nil {
				return nil, err
			}
			s.files[fi.Name()] = string(fd)
		}
	}
	return s, nil
}

// Get returns a loaded file's contents and existence of the file by filename.
func (p *Purse) Get(filename string) (v string, ok bool) {
	p.mux.RLock()
	v, ok = p.files[filename]
	p.mux.RUnlock()
	return
}
