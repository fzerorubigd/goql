package goql

import (
	"sync"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/fzerorubigd/goql/structures"
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
	structures.RegisterTable("types", &typeProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	})

	structures.RegisterField("types", "name", genericName{})
	structures.RegisterField("types", "pkg_name", genericPackageName{})
	structures.RegisterField("types", "pkg_path", genericPackagePath{})
	structures.RegisterField("types", "file", genericFileName{})
	structures.RegisterField("types", "exported", genericIsExported{})
	structures.RegisterField("types", "docs", genericDoc{})
}

func init() {
	registerTypes()
}
