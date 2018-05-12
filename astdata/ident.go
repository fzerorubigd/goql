package astdata

import (
	"go/ast"
)

// IdentType is the normal type name
type IdentType struct {
	embededData

	ident string
}

func (i *IdentType) String() string {
	return i.ident
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
		embededData: embededData{
			pkg:  p,
			fl:   f,
			node: t,
		},
		ident: ident,
	}
}

func getBasicIdent(t string) Definition {
	return &IdentType{
		ident: t,
	}
}
