package astdata

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

// Field is a single field of a structure, a variable, with tag
type Field struct {
	Name string
	Type Definition
	Docs Docs
	Tags reflect.StructTag
}

// Embed is the embedded type in the struct or interface
type Embed struct {
	Definition
	Docs Docs
	Tags reflect.StructTag
}

// Embeds is a list of embedded items
type Embeds []*Embed

// Fields a list of fields
type Fields []*Field

// StructType is the structures in golang source code
type StructType struct {
	pkg    *Package
	Fields Fields
	Embeds Embeds
}

// String convert struct to string
func (s *StructType) String() string {
	if len(s.Embeds) == 0 && len(s.Fields) == 0 {
		return "struct{}"
	}
	res := "struct {\n"
	for e := range s.Embeds {
		res += "\t" + s.Embeds[e].String() + "\n"
	}

	for f := range s.Fields {
		tags := strings.Trim(string(s.Fields[f].Tags), "`")
		if tags != "" {
			tags = "`" + tags + "`"
		}
		res += fmt.Sprintf("\t%s %s %s\n", s.Fields[f].Name, s.Fields[f].Type.String(), tags)
	}
	return res + "}"
}

// Package return the package of this struct
func (s *StructType) Package() *Package {
	return s.pkg
}

func getStruct(p *Package, f *File, t *ast.StructType) Definition {
	res := &StructType{
		pkg: p,
	}
	for _, s := range t.Fields.List {
		if s.Names != nil {
			for i := range s.Names {

				f := Field{
					Name: nameFromIdent(s.Names[i]),
					Type: newType(p, f, s.Type),
				}
				if s.Tag != nil {
					f.Tags = reflect.StructTag(s.Tag.Value)
					f.Tags = f.Tags[1 : len(f.Tags)-1]
				}
				f.Docs = docsFromNodeDoc(s.Doc)
				res.Fields = append(res.Fields, &f)
			}
		} else {
			e := Embed{
				Definition: newType(p, f, s.Type),
			}
			if s.Tag != nil {
				e.Tags = reflect.StructTag(s.Tag.Value)
				e.Tags = e.Tags[1 : len(e.Tags)-1]
			}
			e.Docs = docsFromNodeDoc(s.Doc)
			res.Embeds = append(res.Embeds, &e)
		}
	}

	return res
}
