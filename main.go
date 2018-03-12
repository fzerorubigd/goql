package main

import (
	"github.com/fzerorubigd/goql/internal/parse"
)

const (
	tt int = iota
)

func main() {
	parse.AST(`select *,aa,ss,"swss" ."sw" from xx`)
}
