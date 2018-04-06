# GoQL 
A query language, over Go code, in Go!

[![Build Status](https://travis-ci.org/fzerorubigd/goql.svg)](https://travis-ci.org/fzerorubigd/goql)
[![Coverage Status](https://coveralls.io/repos/github/fzerorubigd/goql/badge.svg?branch=master)](https://coveralls.io/github/fzerorubigd/goql?branch=master)
[![GoDoc](https://godoc.org/github.com/fzerorubigd/goql?status.svg)](https://godoc.org/github.com/fzerorubigd/goql)
[![Go Report Card](https://goreportcard.com/badge/github.com/fzerorubigd/goql/die-github-cache-die)](https://goreportcard.com/report/github.com/fzerorubigd/goql)

*This package is under heavy development, anything may change!*

This is a golang sql driver, to interact with Go code. currently only select is possible, but the insert/update/delete is in todo list.

## Usage 

like any other sql driver in golang, just import the goql package in your code : 

```go
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/fzerorubigd/goql"
)

func main() {
	// open the net/http package
	con, err := sql.Open("goql", "net/http")
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()

	rows, err := con.Query("SELECT name, receiver, def FROM funcs")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var (
			name string
			rec  sql.NullString
			def  string
		)
		rows.Scan(&name, &rec, &def)
		if rec.Valid {
			name = rec.String + "." + name
		}
		fmt.Printf("\nfunc %s , definition : %s", name, def)
	}

}
```

Also there is an example command line is available for more advanced usage in `cmd/goql` by running `go get -u github.com/fzerorubigd/goql/...` the binary is available in your `GOBIN` directory. you can run query against any installed package in your `GOPATH` via this tool.

List of supported tables and fields are available in [docs/table](docs/tables.md)

there is one special type called `definition`. this type is printed as string, but one can use functions to handle special queries. list of supported functions are available at [docs/functions](docs/functions.md) 

also its possible to add new tables/fields/functions using plugins. an example plugin is available at [plugin/goqlimport](plugin/goqlimport/reg_import.go)

currently only supported query is `select` , with `where`,`order` and `limit` some example query : 


```sql
select * from files where the docs is not null
select * from funcs where def = 'func()' and exported
select * from consts order by name desc limit 10, 10
select * from vars where is_struct(def) and name like 's%'
select * from types where is_map(def) and map_key(def) = 'string'
select * from imports where canonical = 'ctx'
```

## Demo 

[![asciicast](https://asciinema.org/a/170483.png)](https://asciinema.org/a/170483)

## TODO

- Write (more) documentation
- UPDATE/INSERT/DELETE support (Yes, code generation with sql :) )
