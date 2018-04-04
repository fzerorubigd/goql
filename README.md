# GoQL 
A query language, over Go code, in Go!

[![Build Status](https://travis-ci.org/fzerorubigd/goql.svg)](https://travis-ci.org/fzerorubigd/goql)
[![Coverage Status](https://coveralls.io/repos/github/fzerorubigd/goql/badge.svg?branch=master)](https://coveralls.io/github/fzerorubigd/goql?branch=master)
[![GoDoc](https://godoc.org/github.com/fzerorubigd/goql?status.svg)](https://godoc.org/github.com/fzerorubigd/goql)
[![Go Report Card](https://goreportcard.com/badge/github.com/fzerorubigd/goql/die-github-cache-die)](https://goreportcard.com/report/github.com/fzerorubigd/goql)

*This package is under heavy development, anything may change!*

This is a subset of sql, over Golang code. the idea is to interact with Go code in sql. the tables are dynamic and adding column/table is possible.

## What is this?

```
go get -u github.com/fzerorubigd/goql/...

```

A test command line is built in your GOBIN directory 

```
goql --package="fmt" "select * from file"
goql --package="fmt" "select * from funcs"
goql --package="fmt" "select * from consts"
goql --package="fmt" "select * from vars"
goql --package="fmt" "select * from types"
goql --package="fmt" "select * from imports"
```

also some operators are available: 

```
goql --package="fmt" "select name from funcs where receiver is not null and name like '%print' order by name desc limit 10,1"
```

## Demo 

[![asciicast](https://asciinema.org/a/170483.png)](https://asciinema.org/a/170483)

## TODO

its in alpha stage, there is a long todo list :

- Write documentation
- Definition type and operator 
- UPDATE/INSERT/DELETE support (Yes, code generation with sql)
