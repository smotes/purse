# purse [![Build Status](https://drone.io/github.com/smotes/purse/status.png)](https://drone.io/github.com/smotes/purse/latest) [![godoc](https://img.shields.io/badge/godoc-documentation-blue.svg)](http://godoc.org/github.com/smotes/purse)

purse is a persistence-layer Go package for loading SQL file contents for use in Go programs.

**Disclaimer**: *purse is not a query builder or ORM package, but rather a way to organize and load SQL files*.

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

Load the `./sql` directory using `purse.Load()` to have access to the file's contents.

```go
// Load all SQL files from specified directory into a map
ps, err := purse.Load(filepath.Join(".", "sql"))

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
