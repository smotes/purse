package main // import "github.com/smotes/purse/cmd/purse"

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	flag.StringVar(&out, "out", "./", "directory of the output source file")
	flag.StringVar(&file, "file", "out.go", "name of the output source file")
	flag.StringVar(&name, "name", "gen", "variable name of the generated Purse struct")
	flag.StringVar(&pack, "pack", "", "name of the go package for the generated source file")
	flag.Parse()
}

func main() {
	validate(in, errors.New("must provide directory of input SQL file(s)"))
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

	ctx := &context{
		Varname: name,
		Package: pack,
		Files:   data,
	}

	cntnts := contentsHead + contentsBodyStruct + "\n" + contentsBodyVar

	if out != "./" {
		ctx.Varname = strings.Title(name)
		cntnts = contentsHead + contentsBodyVar

		tmplCommon, err := template.New(name).Parse(
				contentsHead + contentsBodyStruct)
		handle(err)

		fCommon, err := os.Create(filepath.Join(out, pack+".go"))
		handle(err)

		err = tmplCommon.Execute(fCommon, ctx)
		handle(err)
	}

	tmpl, err := template.New(name).Parse(cntnts)
	handle(err)

	f, err := os.Create(filepath.Join(out, file))
	handle(err)

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
	contentsHead = `package {{.Package}}

`

	contentsBodyVar = `// {{.Varname}} is a *GenPurse.
	var {{.Varname}} = &GenPurse{
	files: map[string]string{
		{{range $key, $val := .Files}}
			"{{$key}}": {{$val}},
		{{end}}
	},
}
`

	contentsBodyStruct = `import (
	"sync"
)

// GenPurse is a literal implementation of a Purse that is programmatically
// generated from SQL file contents within a directory via go generate.
type GenPurse struct {
	mu sync.RWMutex
	files map[string]string
}

// Get takes a filename and returns a query if it is found within the relevant
// map.  If filename is not found, ok will return false.
func (p *GenPurse) Get(filename string) (v string, ok bool) {
	p.mu.RLock()
	v, ok = p.files[filename]
	p.mu.RUnlock()
	return
}
`
)
