// +build go1.8

package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strings"

	_ "github.com/fzerorubigd/goql"
	"github.com/fzerorubigd/goql/plugin/goqlimport"
	"github.com/ogier/pflag"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
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
	rows, err := c.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	columns, err := rows.ColumnTypes()
	if err != nil {
		log.Fatal(err)
	}

	// Scan needs an array of pointers to the values it is setting
	// This creates the object and sets the values correctly
	var objects []map[int]interface{}
	result := make([][]string, 0)
	for rows.Next() {
		values := make([]interface{}, len(columns))
		object := map[int]interface{}{}
		for i, column := range columns {
			object[i] = reflect.New(column.ScanType()).Interface()
			values[i] = object[i]
		}
		err = rows.Scan(values...)
		if err != nil {
			log.Fatal(err)
		}
		stringData := make([]string, len(cols))
		for i, raw := range object {
			if raw == nil {
				stringData[i] = "<nil>"
			} else {
				stringData[i] = cast.ToString(raw)
			}
		}
		result = append(result, stringData)
	}

	for i, object := range objects {
		stringData := make([]string, len(cols))
		for i, raw := range object {
			if raw == nil {
				stringData[i] = "<nil>"
			} else {
				stringData[i] = cast.ToString(raw)
			}
		}
		result[i] = stringData
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
