package astdata

import (
	"go/ast"
)

// IdentType is the normal type name
type IdentType struct {
	pkg *Package
	fl  *File

	ident string
}

func (i *IdentType) String() string {
	return i.ident
}

// Package get the package of ident
func (i *IdentType) Package() *Package {
	return i.pkg
}

// File return the file of the type
func (i *IdentType) File() *File {
	return i.fl
}

// Ident is the ident of this type
func (i *IdentType) Ident() string {
	return i.ident
}

// Compare try to compare this to def
func (i *IdentType) Compare(def Definition) bool {
	return i.String() == def.String()
}

func getIdent(p *Package, f *File, t *ast.Ident) Definition {
	// ident is the simplest one (I was wrong :)) ).
	ident := nameFromIdent(t)
	//	if isBuiltinIdent(ident) {
	//		p = nil
	//	}
	return &IdentType{
		pkg:   p,
		fl:    f,
		ident: ident,
	}
}

func getBasicIdent(t string) Definition {
	return &IdentType{
		ident: t,
	}
}
