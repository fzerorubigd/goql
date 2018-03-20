package astdata

import (
	"fmt"
	"go/ast"
)

// MapType is the map in the go code
type MapType struct {
	pkg   *Package
	Key   Definition
	Value Definition
}

func (m *MapType) String() string {
	return fmt.Sprintf("map[%s]%s", m.Key.String(), m.Value.String())
}

// Package return the map package
func (m *MapType) Package() *Package {
	return m.pkg
}

func getMap(p *Package, f *File, t *ast.MapType) Definition {
	return &MapType{
		pkg:   p,
		Key:   newType(p, f, t.Key),
		Value: newType(p, f, t.Value),
	}
}
