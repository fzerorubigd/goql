package astdata

import (
	"go/ast"
)

// Function is functions with name, not the func types
type Function struct {
	pkg  *Package
	file *File
	fn   *ast.FuncDecl

	name     string
	docs     Docs
	receiver *Variable
}

// Name return the name of the function
func (f *Function) Name() string {
	return f.name
}

// Docs get the function docs +lazy
func (f *Function) Docs() Docs {
	if f.docs == nil {
		f.docs = docsFromNodeDoc(f.fn.Doc)
	}
	return f.docs
}

// Receiver get the receiver of the function if the function is a method, if not,
// return nil
func (f *Function) Receiver() *Variable {
	return f.receiver
}

// Package get the package of function
func (f *Function) Package() *Package {
	return f.pkg
}

// File get the file of function
func (f *Function) File() *File {
	return f.file
}

// newFunction return a single function annotation
func newFunction(p *Package, fl *File, f *ast.FuncDecl) *Function {
	res := &Function{
		pkg:  p,
		file: fl,
		fn:   f,
		name: nameFromIdent(f.Name),
	}

	if res.fn.Recv != nil {
		n := ""
		if res.fn.Recv.List[0].Names != nil {
			n = nameFromIdent(res.fn.Recv.List[0].Names[0])
		}
		res.receiver = newVariableFromExpr(res.pkg, res.file, n, res.fn.Recv.List[0].Type)

		// a hack for function name
		// TODO : after handling the definition its very simple to use that part
		switch t := res.fn.Recv.List[0].Type.(type) {
		case *ast.Ident:
			res.name = nameFromIdent(t) + "." + res.name
		case *ast.StarExpr:
			res.name = nameFromIdent(t.X.(*ast.Ident)) + "." + res.name
		}
	}
	return res
}
