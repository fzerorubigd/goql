package astdata

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

// File is a single file in a structure
type File struct {
	fileName    string
	packageName string
	pkg         *Package

	docs      Docs
	imports   []*Import
	variables []*Variable
	functions []*Function
	constants []*Constant
	types     []*Type
}

type walker struct {
	src     string
	File    *File
	Package *Package
}

func (fv *walker) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		switch t := node.(type) {
		case *ast.File:
			fv.File.packageName = nameFromIdent(t.Name)
			fv.File.docs = docsFromNodeDoc(t.Doc)
		case *ast.FuncDecl:
			fv.File.functions = append(fv.File.functions, newFunction(fv.Package, fv.File, t))
			return nil // Do not go deeper
		case *ast.GenDecl:
			for i := range t.Specs {
				switch decl := t.Specs[i].(type) {
				case *ast.ImportSpec:
					fv.File.imports = append(fv.File.imports, newImport(fv.Package, fv.File, decl, t.Doc))
				case *ast.ValueSpec:
					if t.Tok.String() == "var" {
						fv.File.variables = append(fv.File.variables, newVariable(fv.Package, fv.File, decl, t.Doc)...)
					} else if t.Tok.String() == "const" {
						var last *Constant
						if len(fv.File.constants) > 0 {
							last = fv.File.constants[len(fv.File.constants)-1]
						}
						fv.File.constants = append(fv.File.constants, newConstant(fv.Package, fv.File, decl, t.Doc, last)...)
					}
				case *ast.TypeSpec:
					fv.File.types = append(fv.File.types, newTypeName(fv.Package, fv.File, decl, t.Doc))
				}
			}
			return nil
		default:
			//fmt.Printf("\n%T=====>%+v", t, t)
		}
	}
	return fv
}

// Package return the package of the file
func (f *File) Package() *Package {
	return f.pkg
}

// FileName return the file name, no path
func (f *File) FileName() string {
	return filepath.Base(f.fileName)
}

// ParseFile try to parse a single file for its annotations
func ParseFile(src string, p *Package) (*File, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	fw := &walker{}
	fw.src = src
	fw.File = &File{}
	fw.Package = p

	ast.Walk(fw, f)
	fw.File.pkg = p
	return fw.File, nil
}
