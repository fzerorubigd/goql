package astdata

import (
	"go/ast"
	"strings"
)

// FuncType is the single function
type FuncType struct {
	pkg *Package

	Parameters []*Variable
	Results    []*Variable
}

func (f *FuncType) getDefinitionWithName(name string) string {
	return name + f.Sign()
}

// Sign return the function sign
func (f *FuncType) Sign() string {
	var args, res []string
	for a := range f.Parameters {
		args = append(args, f.Parameters[a].Type.String())
	}

	for a := range f.Results {
		res = append(res, f.Results[a].Type.String())
	}

	result := "(" + strings.Join(args, ",") + ")"
	if len(res) > 1 {
		result += " (" + strings.Join(res, ",") + ")"
	} else {
		result += " " + strings.Join(res, ",")
	}

	return result
}

// String is the string representation of func type
func (f *FuncType) String() string {
	return "func " + f.Sign()
}

// Package is the func package
func (f *FuncType) Package() *Package {
	return f.pkg
}

func getVariableList(p *Package, fl *File, f *ast.FieldList) []*Variable {
	var res []*Variable
	if f == nil {
		return res
	}
	for i := range f.List {
		n := f.List[i]
		if n.Names != nil {
			for in := range n.Names {
				p := newVariableFromExpr(p, fl, nameFromIdent(n.Names[in]), f.List[i].Type)
				res = append(res, p)
			}
		} else {
			// Its probably without name part (ie return variable)
			p := newVariableFromExpr(p, fl, "", f.List[i].Type)
			res = append(res, p)
		}
	}

	return res
}

func getFunc(p *Package, f *File, t *ast.FuncType) Definition {
	return &FuncType{
		pkg:        p,
		Parameters: getVariableList(p, f, t.Params),
		Results:    getVariableList(p, f, t.Results),
	}
}
