package astdata

import (
	"fmt"
	"go/ast"
	"strconv"
)

// ArrayType is the base array
type ArrayType struct {
	pkg   *Package
	Slice bool
	Len   int
	Type  Definition
}

// EllipsisType is slice type but with ...type definition
type EllipsisType struct {
	*ArrayType
}

// String represent array in string
func (a *ArrayType) String() string {
	if a.Slice {
		return "[]" + a.Type.String()
	}
	return fmt.Sprintf("[%d]%s", a.Len, a.Type.String())
}

// Package return the array package
func (a *ArrayType) Package() *Package {
	return a.pkg
}

// String represent ellipsis array in string
func (e *EllipsisType) String() string {
	return fmt.Sprintf("[...]%s{}", e.Type.String())
}

func getArray(p *Package, f *File, t *ast.ArrayType) Definition {
	slice := t.Len == nil
	ellipsis := false
	l := 0
	if !slice {
		var (
			ls string
		)
		switch t.Len.(type) {
		case *ast.BasicLit:
			ls = t.Len.(*ast.BasicLit).Value
		case *ast.Ellipsis:
			ls = "0"
			ellipsis = true
		}
		l, _ = strconv.Atoi(ls)
	}
	var at Definition = &ArrayType{
		pkg:   p,
		Slice: t.Len == nil,
		Len:   l,
		Type:  newType(p, f, t.Elt),
	}
	if ellipsis {
		at = &EllipsisType{ArrayType: at.(*ArrayType)}
	}
	return at
}
