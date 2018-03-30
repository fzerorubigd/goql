package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strings"

	_ "github.com/fzerorubigd/goql"
	"github.com/fzerorubigd/goql/plugin/goqlimport"
	"github.com/ogier/pflag"
	"github.com/olekukonko/tablewriter"
)

var (
	pkg    = pflag.StringP("package", "p", "net/http", "the package to query against")
	format = pflag.StringP("format", "f", "table", "format of output, json and table ")
	imprt  = pflag.BoolP("go-import", "i", true, "add goimports field to file table? its slower than other fields")
)

func tableWriter(row []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(row)
	for i := range data {
		table.Append(data[i])
	}
	table.Render()
}

func jsonWriter(row []string, data [][]string) {
	l := make(map[string]interface{})
	d := json.NewEncoder(os.Stdout)
	d.SetIndent("", "  ")
	for i := range data {
		for j := range row {
			l[row[j]] = data[i][j]
		}
		err := d.Encode(l)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	pflag.Parse()
	if *imprt {
		goqlimport.Register()
	}

	query := strings.Join(pflag.Args(), " ")
	c, err := sql.Open("goql", *pkg)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	row, err := c.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	cols, err := row.Columns()
	if err != nil {
		log.Fatal(err)
	}

	// Result is your slice string.
	rawResult := make([][]byte, len(cols))
	result := make([][]string, 0)

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	cur := 0
	for row.Next() {
		err = row.Scan(dest...)
		if err != nil {
			log.Fatal(err)
		}
		rr := make([]string, len(cols))
		for i, raw := range rawResult {
			if raw == nil {
				rr[i] = "<nil>"
			} else {
				rr[i] = string(raw)
			}
		}
		result = append(result, rr)
		cur++
	}

	switch strings.ToLower(*format) {
	case "json":
		jsonWriter(cols, result)
	case "table":
		tableWriter(cols, result)
	default:
		log.Fatalf("invalid format %s", *format)
	}

}
