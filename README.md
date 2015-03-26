# purse [![Build Status](https://drone.io/github.com/smotes/purse/status.png)](https://drone.io/github.com/smotes/purse/latest) [![GoDoc](https://godoc.org/github.com/smotes/purse?status.svg)](https://godoc.org/github.com/smotes/purse)

purse is a persistence-layer Go package for loading and embedding SQL file contents into Go programs.

**Disclaimer**: *purse is not a query builder or ORM package, but rather a package and tool used to load and embed SQL files*.

### Rationale

Writing SQL statements using Go strings can quickly become cumbersome and hard to maintain; lacking helpful formatting, syntax highlighting and availability of DBMS-specific tooling such as SQL editors and query planners.

Alternate solutions such as query builders and ORMs are non-portable and require learning additional (often non-idiomatic) Go syntax.

The solution is simple: **keep your SQL inside of SQL files**.

### Setup

First, simply get the package with `go get`:

```bash
$ go get github.com/smotes/purse
```

You can then import it into your Go source file(s):

```go
import (
    ...
    "github.com/smotes/purse"
)
```


### Example

This example assumes there exists a file `query_all.sql` in the `./sql` directory:

```sql
SELECT id, slug, title, created, markdown, html
FROM post
```

Load the `./sql` directory using `purse.New()` to have access to the file's contents.

```go
// Load all SQL files from specified directory into a map
ps, err := purse.New(filepath.Join(".", "sql"))

// Get a file's contents
contents, ok := ps.Get("query_all.sql")
if !ok {
    fmt.Println("SQL file not loaded")
}

// Open the database handler
db, err := sql.Open("postgres", "...")

// Query directly via the database handler
rows, err := db.Query(contents)

// Prepare statements via the database handler
stmt, err := db.Prepare(contents)
```

**Note**: purse only loads files with the `.sql` extension. All other file types in the loaded directory will be ignored.

### Tool

`purse` is also a command line tool to automate the creation of a Purse implementation given a specified
directory of SQL files. Given the directory of SQL files and the directory of the output file,
`purse` will create a Go source file (named `out.go` by default) which contains an implementation and
instantiation (bound to variable gen by default) of a Purse interface driven by a `map[string]string`
literal representing the SQL directory's files' contents.

**Note**: The purse package is meant to be used during development where SQL files are changing often
and need to be reloaded into memory on each program execution. *Contrastingly*, the `purse` tool
is meant to be used in production environments where the SQL files' contents can be embedded
directly into the compiled binary, easing deployment.

To install the tool, use go install:

```bash
$ go get github.com/smotes/purse/cmd/purse
$ go install github.com/smotes/purse/cmd/purse
```

The command syntax is `purse -in="input/dir" -out="output/dir" [-file="out.go"] [-name="gen"]`.

The input directory and output directory paths must either be absolute or relative to the
package using it via go generate, or relative to the current working directory where the
the command was executed.

To override the default output source file name (out.go), provide the optional -file flag.

To override the default variable name (gen) of the generated Purse, provide the optional -name flag.

This process should generally be handled using go generate. Add a comment in one of your go source files, like so:

```go
//go:generate purse -in="./fixtures" -out="."
```

Then run go generate:

```bash
$ go generate
```

And that's it!

**Note**: The `purse` tool depends on certain environment variables to be set to execute properly, namely the
`$GOPACKAGE` variable set automatically when running the go generate command. If you wish to explicitly
use this tool without go generate, you will have to set the output source file's package name by setting this
environment variable.
