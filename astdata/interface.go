package astdata

import "go/ast"

// InterfaceType is the interface in go code
type InterfaceType struct {
	pkg       *Package
	Functions []*Function
	Embeds    []Definition // IdentType or SelectorType
}

func (i *InterfaceType) String() string {
	if len(i.Embeds) == 0 && len(i.Functions) == 0 {
		return "interface{}"
	}

	res := "interface{\n"
	for e := range i.Embeds {
		res += "\t" + i.Embeds[e].String() + "\n"
	}
	for f := range i.Functions {
		res += "\t" + i.Functions[f].Type.getDefinitionWithName(i.Functions[f].name) + "\n"
	}
	return res + "}"
}

// Package get the interface package
func (i *InterfaceType) Package() *Package {
	return i.pkg
}

func getInterface(p *Package, f *File, t *ast.InterfaceType) Definition {
	// TODO : interface may refer to itself I need more time to implement this
	iface := &InterfaceType{}
	for i := range t.Methods.List {
		res := Function{}
		// The method name is mandatory and always 1
		if len(t.Methods.List[i].Names) > 0 {
			res.name = nameFromIdent(t.Methods.List[i].Names[0])

			res.docs = docsFromNodeDoc(t.Methods.List[i].Doc)
			typ := newType(p, f, t.Methods.List[i].Type)
			res.Type = typ.(*FuncType)
			iface.Functions = append(iface.Functions, &res)
		} else {
			// This is the embedded interface
			embed := newType(p, f, t.Methods.List[i].Type)
			iface.Embeds = append(iface.Embeds, embed)
		}

	}
	return iface
}
