package astdata

import (
	"fmt"
	"go/ast"
	"strconv"
)

// ArrayType is the base array
type ArrayType struct {
	pkg   *Package
	fl    *File
	slice bool
	len   int
	def   Definition
}

// EllipsisType is slice type but with ...type definition
type EllipsisType struct {
	*ArrayType
}

// String represent array in string
func (a *ArrayType) String() string {
	if a.slice {
		return "[]" + a.def.String()
	}
	return fmt.Sprintf("[%d]%s", a.len, a.def.String())
}

// Package return the array package
func (a *ArrayType) Package() *Package {
	return a.pkg
}

// File return the file of this type
func (a *ArrayType) File() *File {
	return a.fl
}

// ValueDefinition return the definition of value
func (a *ArrayType) ValueDefinition() Definition {
	return a.def
}

// Len is the len of this array
func (a *ArrayType) Len() int {
	return a.len
}

// Slice means this array is an slice
func (a *ArrayType) Slice() bool {
	return a.slice
}

// String represent ellipsis array in string
func (e *EllipsisType) String() string {
	return fmt.Sprintf("[...]%s{}", e.def.String())
}

// Compare try to compare this to def
func (a *ArrayType) Compare(def Definition) bool {
	return a.String() == def.String()
}

// Compare try to compare this to def
func (e *EllipsisType) Compare(def Definition) bool {
	return e.String() == def.String()
}

func getArray(p *Package, f *File, t *ast.ArrayType) Definition {
	slice := t.Len == nil
	ellipsis := false
	var l int64
	if !slice {
		var (
			ls string
		)
		switch t.Len.(type) {
		case *ast.BasicLit:
			ls = t.Len.(*ast.BasicLit).Value
		case *ast.Ellipsis:
			ls = "0" // TODO : Detect this type size
			ellipsis = true
		}
		l, _ = strconv.ParseInt(ls, 10, 0)
	}
	var at Definition = &ArrayType{
		pkg:   p,
		fl:    f,
		slice: t.Len == nil,
		len:   int(l),
		def:   newType(p, f, t.Elt),
	}
	if ellipsis {
		at = &EllipsisType{ArrayType: at.(*ArrayType)}
	}
	return at
}
