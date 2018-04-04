package goql

import (
	"sync"

	"github.com/fzerorubigd/goql/astdata"
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
	RegisterTable("vars", &variableProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	})

	RegisterField("vars", "name", genericName{})
	RegisterField("vars", "pkg_name", genericPackageName{})
	RegisterField("vars", "pkg_path", genericPackagePath{})
	RegisterField("vars", "file", genericFileName{})
	RegisterField("vars", "exported", genericIsExported{})
	RegisterField("vars", "docs", genericDoc{})
	RegisterField("vars", "def", genericDefinition{})
}

func init() {
	registerVariable()
}
