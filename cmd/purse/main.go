package main // import "github.com/smotes/purse/cmd/purse"

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/smotes/purse"
)

const (
	envar = "GOPACKAGE"
)

var (
	in, out, file, name, pack string
)

func init() {
	flag.StringVar(&in, "in", "", "directory of the input SQL file(s)")
	flag.StringVar(&out, "out", "", "directory of the output source file")
	flag.StringVar(&file, "file", "out.go", "name of the output source file")
	flag.StringVar(&name, "name", "gen", "variable name of the generated Purse struct")
	flag.StringVar(&pack, "pack", "", "name of the go package for the generated source file")
	flag.Parse()
}

func main() {
	validate(in, errors.New("must provide directory of input SQL file(s)"))
	validate(out, errors.New("must provide directory of output source file"))
	if pack == "" {
		pack = os.Getenv(envar)
		validate(pack, errors.New("must provide the name of the go package for the generated source file"))
	}

	mp, err := purse.New(in)
	handle(err)

	data := make(map[string]string, len(mp.Files()))
	for _, name := range mp.Files() {
		s, ok := mp.Get(name)
		if !ok {
			log.Printf("Unable to get file %s", name)
			continue
		}
		data[name] = strconv.Quote(s)
	}

	tmpl, err := template.New(name).Parse(contents)
	handle(err)

	f, err := os.Create(filepath.Join(out, file))
	handle(err)

	ctx := &context{
		Varname: name,
		Package: pack,
		Files:   data}
	err = tmpl.Execute(f, ctx)
	handle(err)
}

func validate(s string, err error) {
	if s == "" {
		handle(err)
	}
}

func handle(err error) {
	if err != nil {
		log.Fatalf("%v", err)
	}
}

type context struct {
	Varname string
	Package string
	Files   map[string]string
}

const (
	contents = `package {{.Package}}

import (
	"sync"
)

// GenPurse is an literal implementation of a Purse that is programmatically generated
// from SQL file contents within a directory via go generate.
type GenPurse struct {
	mu sync.RWMutex
	files map[string]string
}

func (p *GenPurse) Get(filename string) (v string, ok bool) {
	p.mu.RLock()
	v, ok = p.files[filename]
	p.mu.RUnlock()
	return
}

var {{.Varname}} = &GenPurse{
	files: map[string]string{
		{{range $key, $val := .Files}}
			"{{$key}}": {{$val}},
		{{end}}
	},
}
`
)
