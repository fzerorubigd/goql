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
	row, data, err := executor.Execute("net/http", `SELECT * FROM vars WHERE file = 'client.go' `)
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
