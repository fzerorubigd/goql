package goql

import (
	"sync"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/fzerorubigd/goql/structures"
)

type variableProvider struct {
	cache map[string][]interface{}
	lock  *sync.Mutex
}

func (v *variableProvider) Provide(in interface{}) []interface{} {
	v.lock.Lock()
	defer v.lock.Unlock()

	p := in.(*astdata.Package)
	if d, ok := v.cache[p.Path()]; ok {
		return d
	}
	va := p.Variables()
	res := make([]interface{}, len(va))
	for i := range va {
		res[i] = va[i]
	}
	v.cache[p.Path()] = res
	return res
}

func registerVariable() {
	structures.RegisterTable("vars", &variableProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	})

	structures.RegisterField("vars", "name", genericName{})
	structures.RegisterField("vars", "pkg_name", genericPackageName{})
	structures.RegisterField("vars", "pkg_path", genericPackagePath{})
	structures.RegisterField("vars", "file", genericFileName{})
	structures.RegisterField("vars", "exported", genericIsExported{})
	structures.RegisterField("vars", "docs", genericDoc{})
}

func init() {
	registerVariable()
}
