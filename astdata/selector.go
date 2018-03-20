package astdata

import (
	"fmt"
	"go/ast"
)

// SelectorType is the type in another package
type SelectorType struct {
	Type Definition

	file     *File
	selector string
	typeName *Type
	imp      *Import
}

func (s *SelectorType) String() string {
	return s.selector + "." + s.Type.String()
}

// Package is the package of selector
func (s *SelectorType) Package() *Package {
	return s.Type.Package()
}

func getSelector(p *Package, f *File, t *ast.SelectorExpr) Definition {
	switch it := t.X.(type) {
	case *ast.Ident:
		res := &SelectorType{
			Type:     getIdent(p, f, t.Sel).(*IdentType),
			selector: nameFromIdent(it),
			file:     f,
		}
		return res
	default:
		panic(fmt.Sprintf("%T is not supported. please report this (with sample code) to add support for it", it))
	}
}
