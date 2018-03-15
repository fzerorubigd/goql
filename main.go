package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fzerorubigd/goql/executor"
	_ "github.com/fzerorubigd/goql/internal/runtime"
	"github.com/ogier/pflag"
)

var (
	pkg = pflag.StringP("package", "p", "net/http", "the package to query against")
)

func main() {
	pflag.Parse()
	sql := strings.Join(pflag.Args(), " ")
	row, data, err := executor.Execute(*pkg, sql)
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
	fmt.Fprintln(w, "")
	w.Flush()
}
