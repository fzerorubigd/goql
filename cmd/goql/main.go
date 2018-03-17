package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/fzerorubigd/goql/executor"
	_ "github.com/fzerorubigd/goql/internal/runtime"
	"github.com/fzerorubigd/goql/structures"
	"github.com/ogier/pflag"
	"github.com/olekukonko/tablewriter"
)

var (
	pkg    = pflag.StringP("package", "p", "net/http", "the package to query against")
	format = pflag.StringP("format", "f", "table", "format of output, json and table ")
)

func formatCol(v []structures.Valuer) []string {
	s := make([]string, len(v))
	for i := range v {
		s[i] = fmt.Sprint(v[i].Value())
	}
	return s
}

func tableWriter(row []string, data [][]structures.Valuer) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(row)
	for i := range data {
		table.Append(formatCol(data[i]))
	}
	table.Render()
}

func jsonWriter(row []string, data [][]structures.Valuer) {
	l := make(map[string]interface{})
	d := json.NewEncoder(os.Stdout)
	d.SetIndent("", "  ")
	for i := range data {
		for j := range row {
			l[row[j]] = data[i][j].Value()
		}
		err := d.Encode(l)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	pflag.Parse()

	sql := strings.Join(pflag.Args(), " ")
	p, err := astdata.ParsePackage(*pkg)
	if err != nil {
		log.Fatal(err)
	}

	row, data, err := executor.Execute(p, sql)
	if err != nil {
		log.Fatal(err)
	}

	switch strings.ToLower(*format) {
	case "json":
		jsonWriter(row, data)
	case "table":
		tableWriter(row, data)
	default:
		log.Fatalf("invalid format %s", *format)
	}
}
