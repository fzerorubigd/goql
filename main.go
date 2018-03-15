package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/fzerorubigd/goql/executor"
	_ "github.com/fzerorubigd/goql/internal/runtime"
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
	row, data, err := executor.Execute("net/http", `SELECT * FROM funcs where receiver is null`)
	if err != nil {
		log.Fatal(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
	fmt.Fprint(w, "\t")

	for i := range row {
		fmt.Fprintf(w, "%s\t", row[i])
	}

	for i := range data {
		fmt.Fprint(w, "\n\t")
		for j := range data[i] {
			fmt.Fprintf(w, "%v\t", data[i][j].Value())
		}
	}
	w.Flush()
}
