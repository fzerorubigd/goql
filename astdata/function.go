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

	receiverClass   string
	receiverPointer bool

	body string
	def  *FuncType
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

// ReceiverType return the type of the receiver
func (f *Function) ReceiverType() string {
	return f.receiverClass
}

// RecieverPointer means this is a pointer method
func (f *Function) RecieverPointer() bool {
	return f.receiverPointer
}

// Body return the function body
func (f *Function) Body() string {
	return f.body
}

// Definition return the definition of the func
func (f *Function) Definition() Definition {
	return f.def
}

// Func return the definition in correct cast. for faster access without a cast
func (f *Function) Func() *FuncType {
	return f.def
}

// newFunction return a single function annotation
func newFunction(p *Package, fl *File, f *ast.FuncDecl) *Function {
	res := &Function{
		pkg:  p,
		file: fl,
		fn:   f,
		name: nameFromIdent(f.Name),
		def:  getFunc(p, fl, f.Type).(*FuncType),
	}

	if res.fn.Recv != nil {
		n := ""
		if res.fn.Recv.List[0].Names != nil {
			n = nameFromIdent(res.fn.Recv.List[0].Names[0])
		}
		res.receiver = newVariableFromExpr(res.pkg, res.file, n, res.fn.Recv.List[0].Type)
		var def Definition
		var def2 *StarType
		def = res.receiver.def
		def2, res.receiverPointer = def.(*StarType)
		if res.receiverPointer {
			def = def2.Target()
		}
		res.receiverClass = def.String()
	}
	return res
}
