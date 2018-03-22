package astdata

import "go/ast"

// StarType is the pointer of a type
type StarType struct {
	def Definition
	pkg *Package
	fl  *File
}

func (s *StarType) String() string {
	return "*" + s.def.String()
}

// Package get the package name
func (s *StarType) Package() *Package {
	return s.pkg
}

// File is the file which this type is defined in it
func (s *StarType) File() *File {
	return s.fl
}

// Target is the target type of this star type
func (s *StarType) Target() Definition {
	return s.def
}

func getStar(p *Package, f *File, t *ast.StarExpr) Definition {
	return &StarType{
		def: newType(p, f, t.X),
		pkg: p,
		fl:  f,
	}
}
