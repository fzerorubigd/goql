package astdata

import (
	"fmt"
	"go/ast"
	"go/token"
)

// Variable is a string represent of a function parameter
type Variable struct {
	pkg *Package
	fl  *File

	name string
	docs Docs

	def Definition

	caller *ast.CallExpr
	index  int
}

// Name return the name of variable
func (v *Variable) Name() string {
	return v.name
}

// Docs return the variable document
func (v *Variable) Docs() Docs {
	return v.docs
}

// Package the package name of this variable
func (v *Variable) Package() *Package {
	return v.pkg
}

// File the filename of this variable
func (v *Variable) File() *File {
	return v.fl
}

// Definition return the variable definition
func (v *Variable) Definition() Definition {
	if v.def == nil && v.caller != nil {
		v.bind()
	}
	return v.def
}

func getBuiltinFunc(name string, v *Variable) Definition {
	fn, err := getBuiltin().FindFunction(name)
	if err == nil {
		// make and new are exceptions,
		if fn.name == "make" {
			return newType(v.pkg, v.fl, v.caller.Args[0])
		}

		if fn.name == "new" {
			tt := newType(v.pkg, v.fl, v.caller.Args[0])
			return &StarType{
				def: tt,
				pkg: v.pkg,
				fl:  v.fl,
			}
		}
	}

	return nil

}

func getNormalFunc(name string, v *Variable) (Definition, error) {
	fn, err := v.pkg.FindFunction(name)
	if err == nil {
		if len(fn.def.results) <= v.index {
			return nil, fmt.Errorf("%d result is available but want the %d", len(fn.def.results), v.index)
		}
		return fn.def.results[v.index].def, nil
	}
	t, err := checkTypeCast(v.pkg, getBuiltin(), v.caller.Args, name)
	if err != nil {
		return nil, err
	}
	return t, err
}

func getForeignType(pkg *Package, pkgName string, fl *File, foreignTyp Definition) Definition {
	star := false
	if sType, ok := foreignTyp.(*StarType); ok {
		foreignTyp = sType.def
		star = true
	}

	switch ft := foreignTyp.(type) {
	case *IdentType:
		// this is a simple hack. if the type is begin with
		// upper case, then its type on that package,
		// else its a global type
		name := ft.ident
		c := name[0]
		if c >= 'A' && c <= 'Z' {
			var res Definition
			res = &SelectorType{
				selector: pkgName,
				fl:       fl,
				pkg:      pkg,
				ident:    name,
			}
			if star {
				res = &StarType{
					def: res,
				}
			}
			return res
		}
		return foreignTyp

	default:
		// the type is foreign to that package too, WTH :)
		return ft
	}

}

func (v *Variable) bind() error {
	if v.caller != nil {
		switch c := v.caller.Fun.(type) {
		case *ast.Ident:
			name := nameFromIdent(c)
			t := getBuiltinFunc(name, v)
			if t != nil {
				v.def = t
				break
			}
			t, err := getNormalFunc(name, v)
			if err != nil {
				return err
			}
			v.def = t

		case *ast.SelectorExpr:
			var pkg string
			switch c.X.(type) {
			case *ast.Ident:
				pkg = nameFromIdent(c.X.(*ast.Ident))
			case *ast.CallExpr: // TODO : Don't know why, no time for check
				break
			}

			typ := nameFromIdent(c.Sel)
			imprt, err := v.pkg.FindImport(pkg)
			if err != nil {
				// TODO : package currently is not capable of parsing build tags. so ignore this :/
				break
			}
			pkgDef, err := ParsePackage(imprt.path)
			if err != nil {
				return err
			}
			var t Definition
			fn, err := pkgDef.FindFunction(typ)
			if err == nil {
				if len(fn.def.results) <= v.index {
					return fmt.Errorf("%d result is available but want %d", len(fn.def.results), v.index)
				}
				t = fn.def.results[v.index].def
			} else {
				t, err = checkTypeCast(pkgDef, builtin, v.caller.Args, typ)
				if err != nil {
					return err
				}
			}

			foreignTyp := t
			v.def = getForeignType(v.pkg, pkg, v.fl, foreignTyp)
		}
	}
	return nil

}

func newVariableFromValue(p *Package, f *File, name string, index int, e []ast.Expr) *Variable {
	var t Definition
	var caller *ast.CallExpr
	var ok bool
	first := e[0]
	// if the caller is a CallExpr, then late bind will take care of it
	if caller, ok = first.(*ast.CallExpr); !ok {
		switch data := e[index].(type) {
		case *ast.CompositeLit:
			t = newType(p, f, data.Type)
		case *ast.BasicLit:
			switch data.Kind {
			case token.INT:
				t = getBasicIdent("int")
			case token.FLOAT:
				t = getBasicIdent("float64")
			case token.IMAG:
				t = getBasicIdent("complex64")
			case token.CHAR:
				t = getBasicIdent("char")
			case token.STRING:
				t = getBasicIdent("string")
			}
			//default:
			//fmt.Printf("var value => %T", e[index])
			//fmt.Printf("%s", src[data.Pos()-1:data.End()-1])
		}
	}
	return &Variable{
		pkg:    p,
		fl:     f,
		name:   name,
		def:    t,
		caller: caller,
		index:  index,
	}
}

// newVariable return an array of variables in the scope
func newVariable(p *Package, f *File, v *ast.ValueSpec, c *ast.CommentGroup) []*Variable {
	var res []*Variable
	for i := range v.Names {
		name := nameFromIdent(v.Names[i])
		var n *Variable
		if v.Type != nil {
			n = newVariableFromExpr(p, f, name, v.Type)
		} else {
			if len(v.Values) != 0 {
				n = newVariableFromValue(p, f, name, i, v.Values)
			}
		}
		n.docs = docsFromNodeDoc(c, v.Doc)
		res = append(res, n)
	}

	return res
}

func newVariableFromExpr(p *Package, f *File, name string, e ast.Expr) *Variable {
	return &Variable{
		pkg:  p,
		fl:   f,
		name: name,
		def:  newType(p, f, e),
	}

}
