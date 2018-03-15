package main

import (
	"github.com/fzerorubigd/goql/executor"
	_ "github.com/fzerorubigd/goql/internal/runtime"
	"github.com/kr/pretty"
)

func main() {
	/*	p, err := astdata.ParsePackage("net/http")
		if err != nil {
			log.Fatal(err)
		}
		ch := make(chan []interface{}, 3)
		err = structures.GetFields(p, "funcs", ch, "name", "pkg_name", "pkg_path", "file", "method")
		if err != nil {
			log.Fatal(err)
		}

		for i := range ch {
			pretty.Print(i)
		}
	*/
	pretty.Print(executor.Execute("net/http", `SELECT name FROM funcs where name = 'main'`))
}
