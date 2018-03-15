package astdata

// TODO : iota support is very limited and bad

import (
	"go/ast"
)

// Constant is a string represent of a function parameter
type Constant struct {
	pkg *Package
	fl  *File

	v *ast.ValueSpec
	c *ast.CommentGroup

	name string
	docs Docs

	value string
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
		n := Constant{
			pkg: p,
			fl:  f,
		}
		name := nameFromIdent(v.Names[i])
		l := ""
		if last != nil {
			l = last.Value()
		}
		n.value = getConstantValue(v.Values, l)
		n.name = name
		n.docs = docsFromNodeDoc(c, v.Doc)
		last = &n
		res = append(res, &n)
	}

	return res
}
