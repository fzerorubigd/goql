package astdata

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

// Field is a single field of a structure, a variable, with tag
type Field struct {
	name string
	def  Definition
	docs Docs
	tags reflect.StructTag
}

// Name is the field name
func (f *Field) Name() string {
	return f.name
}

// Definition of the field
func (f *Field) Definition() Definition {
	return f.def
}

// Docs return the document of the field
func (f *Field) Docs() Docs {
	return f.docs
}

// Tags return the struct tags of the field
func (f *Field) Tags() reflect.StructTag {
	return f.tags
}

// Embed is the embedded type in the struct or interface
type Embed struct {
	Definition
	docs Docs
	tags reflect.StructTag
}

// Docs return the document of the field
func (e *Embed) Docs() Docs {
	return e.docs
}

// Tags return the struct tags of the field
func (e *Embed) Tags() reflect.StructTag {
	return e.tags
}

// Embeds is a list of embedded items
type Embeds []*Embed

// Fields a list of fields
type Fields []*Field

// StructType is the structures in golang source code
type StructType struct {
	pkg    *Package
	fl     *File
	fields Fields
	embeds Embeds
}

// String convert struct to string
func (s *StructType) String() string {
	if len(s.embeds) == 0 && len(s.fields) == 0 {
		return "struct{}"
	}
	res := "struct {\n"
	for e := range s.embeds {
		res += "\t" + s.embeds[e].String() + "\n"
	}

	for f := range s.fields {
		tags := strings.Trim(string(s.fields[f].tags), "`")
		if tags != "" {
			tags = "`" + tags + "`"
		}
		res += fmt.Sprintf("\t%s %s %s\n", s.fields[f].name, s.fields[f].def.String(), tags)
	}
	return res + "}"
}

// Package return the package of this struct
func (s *StructType) Package() *Package {
	return s.pkg
}

// File return the file of type
func (s *StructType) File() *File {
	return s.fl
}

// Fields return struct fields
func (s *StructType) Fields() Fields {
	return s.fields
}

// Embeds is the embed structures
func (s *StructType) Embeds() Embeds {
	return s.embeds
}

func getStruct(p *Package, f *File, t *ast.StructType) Definition {
	res := &StructType{
		pkg: p,
		fl:  f,
	}
	for _, s := range t.Fields.List {
		if s.Names != nil {
			for i := range s.Names {

				f := Field{
					name: nameFromIdent(s.Names[i]),
					def:  newType(p, f, s.Type),
				}
				if s.Tag != nil {
					f.tags = reflect.StructTag(s.Tag.Value)
					f.tags = f.tags[1 : len(f.tags)-1]
				}
				f.docs = docsFromNodeDoc(s.Doc)
				res.fields = append(res.fields, &f)
			}
		} else {
			e := Embed{
				Definition: newType(p, f, s.Type),
			}
			if s.Tag != nil {
				e.tags = reflect.StructTag(s.Tag.Value)
				e.tags = e.tags[1 : len(e.tags)-1]
			}
			e.docs = docsFromNodeDoc(s.Doc)
			res.embeds = append(res.embeds, &e)
		}
	}

	return res
}
