package runtime

import (
	"sync"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/fzerorubigd/goql/structures"
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

func (canonicalCol) Value(in interface{}) structures.String {
	im := in.(*astdata.Import)
	if im.Canonical() == "" {
		return structures.String{Null: true}
	}
	return structures.String{String: im.Canonical()}
}

type pathCol struct{}

func (pathCol) Value(in interface{}) structures.String {
	im := in.(*astdata.Import)
	return structures.String{String: im.Path()}
}

type packageCol struct{}

func (packageCol) Value(in interface{}) structures.String {
	im := in.(*astdata.Import)
	return structures.String{String: im.TargetPackage()}
}

func registerImport() {
	structures.RegisterTable("imports", &importProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	})

	structures.RegisterField("imports", "pkg_name", genericPackageName{})
	structures.RegisterField("imports", "pkg_path", genericPackagePath{})
	structures.RegisterField("imports", "file", genericFileName{})
	structures.RegisterField("imports", "docs", genericDoc{})
	structures.RegisterField("imports", "canonical", canonicalCol{})
	structures.RegisterField("imports", "path", pathCol{})
	structures.RegisterField("imports", "package", packageCol{})
}

func init() {
	registerImport()
}
