// purse is a Go package and tool for loading and embedding SQL file contents into Go programs.
//
// To install the package, use go install:
//
// 		$ go get -u github.com/smotes/purse
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
//		$ go get -u github.com/smotes/purse/cmd/purse
//		$ purse [args]
//
// The command syntax is:
//
//		$ purse -in="input/dir" [-out="output/dir"] [-file="out.go"] [-name="gen"] [-pack="main"]
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
