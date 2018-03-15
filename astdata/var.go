package astdata

import (
	"go/ast"
)

// Variable is a string represent of a function parameter
type Variable struct {
	pkg *Package
	fl  *File

	name string
	docs Docs
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

// newVariable return an array of variables in the scope
func newVariable(p *Package, f *File, v *ast.ValueSpec, c *ast.CommentGroup) []*Variable {
	var res []*Variable
	for i := range v.Names {
		name := nameFromIdent(v.Names[i])
		res = append(res, &Variable{
			pkg:  p,
			fl:   f,
			docs: docsFromNodeDoc(c, v.Doc),
			name: name,
		})

	}

	return res
}

func newVariableFromExpr(p *Package, f *File, name string, e ast.Expr) *Variable {
	return &Variable{
		pkg:  p,
		fl:   f,
		name: name,
	}

}
