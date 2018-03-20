package astdata

import "go/ast"

// IdentType is the normal type name
type IdentType struct {
	pkg   *Package
	Ident string
}

func (i *IdentType) String() string {
	return i.Ident
}

// Package get the package of ident
func (i *IdentType) Package() *Package {
	return i.pkg
}

func getIdent(p *Package, _ *File, t *ast.Ident) Definition {
	// ident is the simplest one (I was wrong :)) ).
	ident := nameFromIdent(t)
	//	if isBuiltinIdent(ident) {
	//		p = nil
	//	}
	return &IdentType{
		pkg:   p,
		Ident: ident,
	}
}

func getBasicIdent(t string) Definition {
	return &IdentType{
		Ident: t,
	}
}
