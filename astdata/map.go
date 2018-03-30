package astdata

import (
	"fmt"
	"go/ast"
)

// MapType is the map in the go code
type MapType struct {
	pkg   *Package
	fl    *File
	key   Definition
	value Definition
}

func (m *MapType) String() string {
	return fmt.Sprintf("map[%s]%s", m.key.String(), m.value.String())
}

// Package return the map package
func (m *MapType) Package() *Package {
	return m.pkg
}

// File is the file of the type
func (m *MapType) File() *File {
	return m.fl
}

// Key type definition
func (m *MapType) Key() Definition {
	return m.key
}

// Value is the definition of the value type
func (m *MapType) Value() Definition {
	return m.value
}

func getMap(p *Package, f *File, t *ast.MapType) Definition {
	return &MapType{
		pkg:   p,
		fl:    f,
		key:   newType(p, f, t.Key),
		value: newType(p, f, t.Value),
	}
}
