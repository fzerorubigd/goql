package astdata

import (
	"os"
	"path/filepath"
)

// Package is the package in one place
type Package struct {
	path string
	name string

	files []*File
}

// Path of the package
func (p *Package) Path() string {
	return p.path
}

// Name of the package
func (p *Package) Name() string {
	return p.name
}

func parsePackageFullPath(path, folder string) (*Package, error) {
	if p := getCache(folder); p != nil {
		return p, nil
	}

	var (
		p = &Package{}
		e error
	)
	p.path = path
	e = filepath.Walk(
		folder,
		func(path string, f os.FileInfo, _ error) error {
			data, err := getGoFileContent(path, folder, f)
			if err != nil || data == "" {
				return err
			}
			fl, err := ParseFile(data, p)
			if err != nil {
				return err
			}
			fl.fileName = path
			p.files = append(p.files, fl)
			if p.name == "" {
				p.name = fl.packageName
			}

			return nil
		},
	)
	if e != nil {
		return nil, e
	}
	setCache(folder, p)

	return p, nil
}

// ParsePackage is here for loading a single package and parse all files in it
// if the package is imported from another package, the other parameter is required for
// checking vendors of that package.
func ParsePackage(path string, packages ...string) (*Package, error) {
	folder, err := translateToFullPath(path, packages...)
	if err != nil {
		return nil, err
	}

	return parsePackageFullPath(path, folder)
}
