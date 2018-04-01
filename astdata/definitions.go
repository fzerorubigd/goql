package astdata

import (
	"fmt"
	"go/ast"
	"go/parser"
)

// Definition is the interface for all types without name
type Definition interface {
	fmt.Stringer
	// Package return the package name of the type
	Package() *Package
	// Compare two definition
	// TODO : this is here, but the implementation is not complete. for {Selector,Map,Struct,Interface,Func}Type we need to check for selector type canonical name and pkg name
	Compare(Definition) bool
}

func newType(p *Package, f *File, e ast.Expr) Definition {
	switch t := e.(type) {
	case *ast.Ident:
		return getIdent(p, f, t)
	case *ast.StarExpr:
		return getStar(p, f, t)
	case *ast.ArrayType:
		return getArray(p, f, t)
	case *ast.MapType:
		return getMap(p, f, t)
	case *ast.StructType:
		return getStruct(p, f, t)
	case *ast.SelectorExpr:
		return getSelector(p, f, t)
	case *ast.ChanType:
		return getChannel(p, f, t)
	case *ast.FuncType:
		return getFunc(p, f, t)
	case *ast.InterfaceType:
		return getInterface(p, f, t)
	case *ast.CompositeLit:
		if tmp, ok := t.Type.(*ast.ArrayType); ok {
			return getArray(p, f, tmp)
		}
	}
	return nil
}

// NewDefinition try to extract definition from string
func NewDefinition(s string) (Definition, error) {
	d, err := parser.ParseExpr(s)
	if err != nil {
		return nil, err
	}
	t := newType(&Package{}, &File{}, d)
	if t == nil {
		return nil, fmt.Errorf("not supported type, maybe not a definition of a type?")
	}
	return t, nil
}
