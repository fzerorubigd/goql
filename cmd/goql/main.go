package main

import (
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
	pkg = pflag.StringP("package", "p", "net/http", "the package to query against")
)

func formatCol(v []structures.Valuer) []string {
	s := make([]string, len(v))
	for i := range v {
		s[i] = fmt.Sprint(v[i].Value())
	}
	return s
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(row)
	for i := range data {
		table.Append(formatCol(data[i]))
	}
	table.Render()
}
