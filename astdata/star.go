package astdata

import "go/ast"

// StarType is the pointer of a type
type StarType struct {
	Target Definition
	pkg    *Package
	file   *File
}

func (s *StarType) String() string {
	return "*" + s.Target.String()
}

// Package get the package name
func (s *StarType) Package() *Package {
	return s.pkg
}

func getStar(p *Package, f *File, t *ast.StarExpr) Definition {
	return &StarType{
		Target: newType(p, f, t.X),
		pkg:    p,
		file:   f,
	}
}
