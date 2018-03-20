package main

import (
	"fmt"
	"log"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/fzerorubigd/goql/executor"
	"github.com/fzerorubigd/goql/internal/parse"
	_ "github.com/fzerorubigd/goql/internal/runtime"
	"github.com/fzerorubigd/goql/structures"
	"github.com/kr/pretty"
)

type fn struct {
}

func (fn) Execute(in ...structures.Valuer) (structures.Valuer, error) {
	s := ""
	for i := range in {
		s += fmt.Sprint(in[i].Value())
	}
	return structures.String{String: s}, nil
}

func main() {
	structures.RegisterFunction("test", fn{})
	ast, err := parse.AST("SELECT test(name, 'ss') FROM funcs where test(name, 'ss') = 'advancess'")
	if err != nil {
		log.Fatal(err)
	}
	p, err := astdata.ParsePackage("fmt")
	if err != nil {
		log.Fatal(err)
	}
	pretty.Print(executor.Execute(p, ast))
}
