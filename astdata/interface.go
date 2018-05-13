package astdata

import (
	"go/ast"
)

// InterfaceType is the interface in go code
type InterfaceType struct {
	embededData

	functions []*Function
	embeds    Embeds
}

func (i *InterfaceType) String() string {
	if len(i.embeds) == 0 && len(i.functions) == 0 {
		return "interface{}"
	}

	res := "interface{\n"
	for e := range i.embeds {
		res += "\t" + i.embeds[e].String() + "\n"
	}
	for f := range i.functions {
		res += "\t" + i.functions[f].def.getDefinitionWithName(i.functions[f].name) + "\n"
	}
	return res + "}"
}

// Functions return the functions in the interface
func (i *InterfaceType) Functions() []*Function {
	return i.functions
}

// Embeds is the embedded interfaces
func (i *InterfaceType) Embeds() Embeds {
	return i.embeds
}

// Compare try to compare this to def
func (i *InterfaceType) Compare(def Definition) bool {
	return i.String() == def.String()
}

func getInterface(p *Package, f *File, t *ast.InterfaceType) Definition {
	// TODO : interface may refer to itself I need more time to implement this
	iface := &InterfaceType{
		embededData: embededData{
			pkg:  p,
			fl:   f,
			node: t,
		},
	}
	for i := range t.Methods.List {
		res := Function{}
		// The method name is mandatory and always 1
		if len(t.Methods.List[i].Names) > 0 {
			res.name = nameFromIdent(t.Methods.List[i].Names[0])

			res.docs = docsFromNodeDoc(t.Methods.List[i].Doc)
			typ := newType(p, f, t.Methods.List[i].Type)
			res.def = typ.(*FuncType)
			iface.functions = append(iface.functions, &res)
		} else {
			// This is the embedded interface
			embed := &Embed{
				def:  newType(p, f, t.Methods.List[i].Type),
				docs: docsFromNodeDoc(t.Methods.List[i].Doc),
				// tags are always empty. (struct tag :)
			}
			iface.embeds = append(iface.embeds, embed)
		}

	}
	return iface
}
