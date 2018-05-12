package astdata

import (
	"go/ast"
)

// SelectorType is the type in another package
type SelectorType struct {
	embededData

	selector string
	resolved string
	ident    string
	imp      *Import
}

func (s *SelectorType) String() string {
	return s.selector + "." + s.ident
}

// Selector is the selector type
func (s *SelectorType) Selector() string {
	return s.selector
}

// Ident is the ident after dot
func (s *SelectorType) Ident() string {
	return s.ident
}

// Import is the import of this selector
func (s *SelectorType) Import() *Import {
	return s.imp
}

// Compare try to compare this to def
func (s *SelectorType) Compare(def Definition) bool {
	// TODO : for this type, structs/map/func/interface that can contain a SelectorType we should check both canonical and real pkg
	return s.String() == def.String()
}

func getSelector(p *Package, f *File, t *ast.SelectorExpr) Definition {
	it := t.X.(*ast.Ident)
	res := &SelectorType{
		embededData: embededData{
			pkg:  p,
			fl:   f,
			node: t,
		},
		ident:    nameFromIdent(t.Sel),
		selector: nameFromIdent(it),
	}

	res.resolved = f.resolvePkg(res.selector)

	for i := range f.imports {
		if f.imports[i].Canonical() == res.selector || f.imports[i].TargetPackage() == res.selector {
			res.imp = f.imports[i]
		}
	}
	return res
}
