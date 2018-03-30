package goql

import (
	"sync"

	"github.com/fzerorubigd/goql/astdata"
)

type typeProvider struct {
	cache map[string][]interface{}
	lock  *sync.Mutex
}

func (t *typeProvider) Provide(in interface{}) []interface{} {
	t.lock.Lock()
	defer t.lock.Unlock()

	p := in.(*astdata.Package)
	if d, ok := t.cache[p.Path()]; ok {
		return d
	}
	fs := p.Types()
	res := make([]interface{}, len(fs))
	for i := range fs {
		res[i] = fs[i]
	}
	t.cache[p.Path()] = res
	return res
}

func registerTypes() {
	RegisterTable("types", &typeProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	})

	RegisterField("types", "name", genericName{})
	RegisterField("types", "pkg_name", genericPackageName{})
	RegisterField("types", "pkg_path", genericPackagePath{})
	RegisterField("types", "file", genericFileName{})
	RegisterField("types", "exported", genericIsExported{})
	RegisterField("types", "docs", genericDoc{})
}

func init() {
	registerTypes()
}
