package astdata

import (
	"fmt"
)

// FindImport try to find an import in a file
func (f *File) FindImport(pkg string) (*Import, error) {
	if pkg == "" || pkg == "_" || pkg == "." {
		return nil, fmt.Errorf("import with path _/. or empty is invalid")
	}
	for i := range f.imports {
		if f.imports[i].TargetPackage() == pkg || f.imports[i].Canonical() == pkg || f.imports[i].Path() == pkg {
			return f.imports[i], nil
		}
	}

	return nil, fmt.Errorf("pkg %s is not found in %s", pkg, f.FileName())
}

// FindImport try to find an import in a package
func (p *Package) FindImport(pkg string) (*Import, error) {
	for i := range p.files {
		if i, err := p.files[i].FindImport(pkg); err == nil {
			return i, nil
		}
	}

	return nil, fmt.Errorf("pkg %s is not found in %s", pkg, p.Name())
}

// FindType try to find type in file
func (f *File) FindType(t string) (*Type, error) {
	for i := range f.types {
		if f.types[i].Name() == t {
			return f.types[i], nil
		}
	}
	return nil, fmt.Errorf("type %s is not found in %s", t, f.FileName())
}

// FindType try to find type in package
func (p *Package) FindType(t string) (*Type, error) {
	for i := range p.files {
		if ty, err := p.files[i].FindType(t); err == nil {
			return ty, nil
		}
	}
	return nil, fmt.Errorf("type %s is not found in %s", t, p.Name())
}

// FindConstant try to find constant in package
func (f *File) FindConstant(t string) (*Constant, error) {
	for i := range f.constants {
		if f.constants[i].Name() == t {
			return f.constants[i], nil
		}
	}
	return nil, fmt.Errorf("const %s is not found in %s", t, f.FileName())
}

// FindConstant try to find constant in package
func (p *Package) FindConstant(t string) (*Constant, error) {
	for i := range p.files {
		if ct, err := p.files[i].FindConstant(t); err == nil {
			return ct, nil
		}
	}
	return nil, fmt.Errorf("const %s is not found in %s", t, p.Name())
}

// FindFunction try to find function in file
func (f *File) FindFunction(t string) (*Function, error) {
	for i := range f.functions {
		if f.functions[i].Receiver() != nil {
			continue
		}
		if f.functions[i].Name() == t {
			return f.functions[i], nil
		}
	}
	return nil, fmt.Errorf("function %s is not found in %s", t, f.FileName())
}

// FindFunction try to find function in package
func (p *Package) FindFunction(t string) (*Function, error) {
	for i := range p.files {
		if ct, err := p.files[i].FindFunction(t); err == nil {
			return ct, nil
		}
	}
	return nil, fmt.Errorf("function %s is not found in %s", t, p.Name())
}

// FindMethod try to find function in file
func (f *File) FindMethod(t string, fn string) (*Function, error) {
	for i := range f.functions {
		if f.functions[i].Receiver() == nil {
			continue
		}
		if f.functions[i].Name() == fn && f.functions[i].ReceiverType() == t {
			return f.functions[i], nil
		}
	}
	return nil, fmt.Errorf("function %s is not found in %s", t, f.FileName())
}

// FindMethod try to find function in package
func (p *Package) FindMethod(t string, fn string) (*Function, error) {
	for i := range p.files {
		if ct, err := p.files[i].FindMethod(t, fn); err == nil {
			return ct, nil
		}
	}
	return nil, fmt.Errorf("function %s is not found in %s", t, p.Name())
}

// FindVariable try to find variable in file
func (f *File) FindVariable(t string) (*Variable, error) {
	for i := range f.variables {
		if f.variables[i].Name() == t {
			return f.variables[i], nil
		}
	}
	return nil, fmt.Errorf("variable %s is not found in %s", t, f.FileName())
}

// FindVariable try to find variable in package
func (p *Package) FindVariable(t string) (*Variable, error) {
	for i := range p.files {
		if ct, err := p.files[i].FindVariable(t); err == nil {
			return ct, nil
		}
	}
	return nil, fmt.Errorf("variable %s is not found in %s", t, p.Name())
}
