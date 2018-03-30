package astdata

import (
	"fmt"
	"go/ast"
)

// Definition is the interface for all types without name
type Definition interface {
	fmt.Stringer
	// Package return the package name of the type
	Package() *Package
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
	default:
		return nil
	}
}
