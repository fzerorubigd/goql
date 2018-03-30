package goql

import (
	"sync"

	"github.com/fzerorubigd/goql/astdata"
)

type importProvider struct {
	cache map[string][]interface{}
	lock  *sync.Mutex
}

func (v *importProvider) Provide(in interface{}) []interface{} {
	v.lock.Lock()
	defer v.lock.Unlock()

	p := in.(*astdata.Package)
	if d, ok := v.cache[p.Path()]; ok {
		return d
	}
	va := p.Imports()
	res := make([]interface{}, len(va))
	for i := range va {
		res[i] = va[i]
	}
	v.cache[p.Path()] = res
	return res
}

type canonicalCol struct{}

func (canonicalCol) Value(in interface{}) String {
	im := in.(*astdata.Import)
	if im.Canonical() == "" {
		return String{Null: true}
	}
	return String{String: im.Canonical()}
}

type pathCol struct{}

func (pathCol) Value(in interface{}) String {
	im := in.(*astdata.Import)
	return String{String: im.Path()}
}

type packageCol struct{}

func (packageCol) Value(in interface{}) String {
	im := in.(*astdata.Import)
	return String{String: im.TargetPackage()}
}

func registerImport() {
	RegisterTable("imports", &importProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	})

	RegisterField("imports", "pkg_name", genericPackageName{})
	RegisterField("imports", "pkg_path", genericPackagePath{})
	RegisterField("imports", "file", genericFileName{})
	RegisterField("imports", "docs", genericDoc{})
	RegisterField("imports", "canonical", canonicalCol{})
	RegisterField("imports", "path", pathCol{})
	RegisterField("imports", "package", packageCol{})
}

func init() {
	registerImport()
}
