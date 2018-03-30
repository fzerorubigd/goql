// Package goqlimport is a prof of concept, just to show how to add a new column to a table
// this is slow, just for demonstration
package goqlimport

import (
	"bytes"
	"fmt"

	"github.com/fzerorubigd/goql"
	"github.com/fzerorubigd/goql/astdata"
	"golang.org/x/tools/imports"
)

type importCheck struct {
	opt *imports.Options
}

func (ic importCheck) Value(in interface{}) goql.Bool {
	// it must be astdata.File
	fl := in.(*astdata.File)
	var src = []byte(fl.Source())
	dst, err := imports.Process(fl.FullPath(), src, ic.opt)
	b := err == nil && bytes.Compare(src, dst) == 0

	fmt.Println(src)
	fmt.Println(dst)
	return goql.Bool{
		Bool: b,
	}
}

func newImportCheck() goql.BoolValuer {
	return importCheck{
		opt: &imports.Options{
			Fragment:   false,
			AllErrors:  false,
			Comments:   true,
			TabIndent:  true,
			TabWidth:   8,
			FormatOnly: true,
		}}
}

// Register the field. Since this is slow one, let the user choose to add it or not
func Register() {
	goql.RegisterField("files", "goimport", newImportCheck())
}
