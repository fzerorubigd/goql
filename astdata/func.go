package astdata

import (
	"go/ast"
	"strings"
)

// FuncType is the single function
type FuncType struct {
	pkg *Package
	fl  *File

	parameters []*Variable
	results    []*Variable
}

func (f *FuncType) getDefinitionWithName(name string) string {
	return name + f.Sign()
}

// Sign return the function sign
func (f *FuncType) Sign() string {
	var args, res []string
	for a := range f.parameters {
		args = append(args, f.parameters[a].def.String())
	}

	for a := range f.results {
		res = append(res, f.results[a].def.String())
	}

	result := "(" + strings.Join(args, ", ") + ")"
	if len(res) > 1 {
		result += " (" + strings.Join(res, ", ") + ")"
	} else {
		result += " " + strings.Join(res, ", ")
	}

	return result
}

// String is the string representation of func type
func (f *FuncType) String() string {
	return f.getDefinitionWithName("func ")
}

// Package is the func package
func (f *FuncType) Package() *Package {
	return f.pkg
}

// File of the type
func (f *FuncType) File() *File {
	return f.fl
}

// Parameters is the parameter of the function
func (f *FuncType) Parameters() []*Variable {
	return f.parameters
}

// Results is the result of the functions
func (f *FuncType) Results() []*Variable {
	return f.results
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
		fl:         f,
		parameters: getVariableList(p, f, t.Params),
		results:    getVariableList(p, f, t.Results),
	}
}
