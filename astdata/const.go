package astdata

// TODO : iota support is very limited and bad

import (
	"go/ast"
	"go/token"
)

// Constant is a string represent of a function parameter
type Constant struct {
	pkg *Package
	fl  *File

	name string
	docs Docs

	value string

	def Definition

	caller *ast.CallExpr
	index  int
}

// Package return the constant package
func (c *Constant) Package() *Package {
	return c.pkg
}

// File return the constant declaration file
func (c *Constant) File() *File {
	return c.fl
}

// Name return the name of constant
func (c *Constant) Name() string {
	return c.name
}

// Docs is the documents of the constant
func (c *Constant) Docs() Docs {
	return c.docs
}

// Value is the value of the constant (not very accurate , the iota is not supported)
func (c *Constant) Value() string {
	return c.value
}

// Definition return the constant definition
func (c *Constant) Definition() Definition {
	return c.def
}

func constantFromValue(p *Package, f *File, name string, indx int, e []ast.Expr) *Constant {
	var t Definition
	var caller *ast.CallExpr
	var ok bool
	if len(e) == 0 {
		return &Constant{
			pkg:  p,
			name: name,
			fl:   f,
		}
	}
	first := e[0]
	if caller, ok = first.(*ast.CallExpr); !ok {
		switch data := e[indx].(type) {
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
				//default:
				//fmt.Printf("var value => %T", e[index])
				//fmt.Printf("%s", src[data.Pos()-1:data.End()-1])
			}
		case *ast.Ident:
			t = getIdent(p, f, data)
		}
	}
	return &Constant{
		pkg:    p,
		fl:     f,
		name:   name,
		def:    t,
		caller: caller,
		index:  indx,
	}
}
func constantFromExpr(p *Package, f *File, name string, e ast.Expr) *Constant {
	return &Constant{
		pkg:  p,
		name: name,
		def:  newType(p, f, e),
		fl:   f,
	}
}

func getConstantValue(a []ast.Expr, lastVal string) string {
	if len(a) == 0 {
		return lastVal
	}
	switch first := a[0].(type) {
	case *ast.BasicLit:
		return first.Value
	default:
		//fmt.Printf("%T ==> %+v", first, first)
		return "NotSupportedYet"
	}
}

// newConstant return an array of constant in the scope
func newConstant(p *Package, f *File, v *ast.ValueSpec, c *ast.CommentGroup, last *Constant) []*Constant {
	var res []*Constant

	for i := range v.Names {
		var n = &Constant{}
		name := nameFromIdent(v.Names[i])

		if v.Type != nil {
			n = constantFromExpr(p, f, name, v.Type)
		} else {
			n = constantFromValue(p, f, name, i, v.Values)
		}

		l := ""
		if last != nil {
			l = last.Value()
		}
		n.value = getConstantValue(v.Values, l)
		n.name = name
		n.docs = docsFromNodeDoc(c, v.Doc)
		n.fl = f
		n.pkg = p
		last = n
		res = append(res, n)
	}

	return res
}
