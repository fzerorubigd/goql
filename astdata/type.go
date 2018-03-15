package astdata

import (
	"go/ast"
)

// Type is the type with name in a package
type Type struct {
	pkg  *Package
	file *File

	t *ast.TypeSpec
	d *ast.CommentGroup

	docs Docs
	name string
}

// Docs return the documents
func (t *Type) Docs() Docs {
	return t.docs
}

// Name is the type name
func (t *Type) Name() string {
	return t.name
}

// Package return the package name of this type name
func (t *Type) Package() *Package {
	return t.pkg
}

// File return the filename of this type name
func (t *Type) File() *File {
	return t.file
}

// newTypeName handle a type with name
func newTypeName(p *Package, f *File, t *ast.TypeSpec, c *ast.CommentGroup) *Type {
	return &Type{
		pkg:  p,
		file: f,
		t:    t,
		d:    c,

		name: nameFromIdent(t.Name),
		docs: docsFromNodeDoc(c, t.Doc),
	}

}
