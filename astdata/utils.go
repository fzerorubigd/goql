package astdata

import (
	"fmt"
	"go/ast"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	cache       = make(map[string]*Package)
	builtinType map[string]bool
	lock        = sync.RWMutex{}
	builtin     *Package
)

func translateToFullPath(path string, packages ...string) (string, error) {
	root := runtime.GOROOT()
	p := os.Getenv("GOPATH")
	if p == "" {
		// TODO : go 1.7 has a default gopath value, must check for other values
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		p = filepath.Join(usr.HomeDir, "go")
	}
	gopath := strings.Split(p, ":")
	gopath = append([]string{root}, gopath...)
	var (
		test string
		r    os.FileInfo
		err  error
	)

	// First try to find it from vendors
	for i := range packages {
		p := filepath.Join(packages[i], "vendor")
		test = filepath.Join(p, path)
		r, err = os.Stat(test)
		if err == nil && r.IsDir() {
			return test, nil
		}
	}

	for i := range gopath {
		test = filepath.Join(gopath[i], "src", path)
		r, err = os.Stat(test)
		if err == nil && r.IsDir() {
			return test, nil
		}
	}
	return "", fmt.Errorf("%s is not found in GOROOT or GOPATH", path)
}

func getGoFileContent(path, folder string, f os.FileInfo) (string, error) {
	if f.IsDir() {
		if path != folder {
			return "", filepath.SkipDir
		}
		return "", nil
	}
	if filepath.Ext(path) != ".go" {
		return "", nil
	}
	// ignore test files (for now?)
	_, filename := filepath.Split(path)
	if len(filename) > 8 && filename[len(filename)-8:] == "_test.go" {
		return "", nil
	}
	r, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = r.Close() }()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func setCache(folder string, p *Package) {
	lock.Lock()
	defer lock.Unlock()

	cache[folder] = p
}

func checkTypeCast(p *Package, bi *Package, args []ast.Expr, name string) (Definition, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("it can not be a typecast : %s", name)
	}

	// if the type is in this package, then simply pass an ident type
	// not the actual type, since the actual type is the definition of
	// the type
	_, err := p.FindType(name)
	if err == nil {
		return &IdentType{pkg: p, ident: name}, nil
	}

	return nil, fmt.Errorf("can not find the call for %s", name)
}

func getCache(folder string) *Package {
	lock.RLock()
	defer lock.RUnlock()

	return cache[folder]
}

func nameFromIdent(i *ast.Ident) (name string) {
	if i != nil {
		name = i.String()
	}
	return
}

func getBuiltin() *Package {
	if builtin == nil {
		var err error
		builtin, err = ParsePackage("builtin")
		if err != nil {
			panic(err)
		}
	}
	return builtin
}

func isBuiltinIdent(ident string) bool {
	if builtinType == nil {
		builtinType = make(map[string]bool)
		b := getBuiltin()
		for f := range b.files {
			for t := range b.files[f].types {
				builtinType[b.files[f].types[t].name] = true
			}
		}
	}

	return builtinType[ident]

}
