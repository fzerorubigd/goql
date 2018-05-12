package astdata

import (
	"go/ast"
)

// StarType is the pointer of a type
type StarType struct {
	embededData

	def Definition
}

func (s *StarType) String() string {
	return "*" + s.def.String()
}

// Target is the target type of this star type
func (s *StarType) Target() Definition {
	return s.def
}

// Compare try to compare this to def
func (s *StarType) Compare(def Definition) bool {
	return s.String() == def.String()
}

func getStar(p *Package, f *File, t *ast.StarExpr) Definition {
	return &StarType{
		embededData: embededData{
			pkg:  p,
			fl:   f,
			node: t,
		},
		def: newType(p, f, t.X),
	}
}
