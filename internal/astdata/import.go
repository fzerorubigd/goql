package astdata

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Import is one imported path
type Import struct {
	targetPkg string
	canonical string
	path      string
	docs      Docs
	i         *ast.ImportSpec
	d         *ast.CommentGroup
	folder    string

	pkg *Package
	fl  *File
}

type importWalker struct {
	pkgName  string
	resolved string
}

func (iw *importWalker) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		switch t := node.(type) {
		case *ast.File:
			iw.pkgName = nameFromIdent(t.Name)
		default:
		}
	}

	return iw
}

func peekPackageName(pkg string, base ...string) (string, string) {
	_, name := filepath.Split(pkg)
	folder, err := translateToFullPath(pkg, base...)
	if err != nil {
		return name, ""
	}
	iw := &importWalker{}
	_ = filepath.Walk(
		folder,
		func(path string, f os.FileInfo, _ error) error {
			// TODO : skip tests and no go files
			data, err := getGoFileContent(path, folder, f)
			if err != nil || data == "" {
				return err
			}
			fset := token.NewFileSet()
			fle, err := parser.ParseFile(fset, "", data, parser.PackageClauseOnly)
			if err != nil {
				return nil // try another file?
			}
			iw.resolved = folder
			ast.Walk(iw, fle)
			// no need to continue
			return filepath.SkipDir
		},
	)
	resolved := ""
	if iw.pkgName != "" {
		name = iw.pkgName
		resolved = iw.resolved

	}
	// can not parse it, use the folder name
	return name, resolved
}

func (i *Import) peek() {
	i.targetPkg, i.folder = peekPackageName(i.path, i.pkg.path)
	if i.canonical == "" {
		i.canonical = i.targetPkg
	}
}

// Package is the package where the import is inside it
func (i *Import) Package() *Package {
	return i.pkg
}

// File return the file where the import is inside it
func (i *Import) File() *File {
	return i.fl
}

// TargetPackage return the target package name. (name after package keyword in package ) +lazy
func (i *Import) TargetPackage() string {
	if i.targetPkg == "" {
		i.peek()
	}

	return i.targetPkg
}

// TargetPath is the target package in full path format
func (i *Import) TargetPath() string {
	return i.path
}

// Folder return the actual folder name of the package +lazy
func (i *Import) Folder() string {
	if i.folder == "" {
		i.peek()
	}

	return i.folder
}

// Canonical is the canonical package name, if the package not imported with another name, the package name and
// canonical name are same +lazy
func (i *Import) Canonical() string {
	if i.canonical == "" {
		i.peek()
	}

	return i.canonical
}

// Docs is the docs for this import path
func (i *Import) Docs() Docs {
	return i.docs
}

// newImport extract a new import entry
func newImport(p *Package, f *File, i *ast.ImportSpec, c *ast.CommentGroup) *Import {
	res := &Import{
		pkg:  p,
		fl:   f,
		path: strings.Trim(i.Path.Value, `"`),
		docs: docsFromNodeDoc(c, i.Doc),
	}
	if i.Name != nil {
		res.canonical = i.Name.String()
	}

	return res
}
