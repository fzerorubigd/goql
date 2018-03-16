package runtime

import (
	"sync"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/fzerorubigd/goql/structures"
)

type constProvider struct {
	cache map[string][]interface{}
	lock  *sync.Mutex
}

func (v *constProvider) Provide(in interface{}) []interface{} {
	v.lock.Lock()
	defer v.lock.Unlock()

	p := in.(*astdata.Package)
	if d, ok := v.cache[p.Path()]; ok {
		return d
	}
	va := p.Constants()
	res := make([]interface{}, len(va))
	for i := range va {
		res[i] = va[i]
	}
	v.cache[p.Path()] = res
	return res
}

func registerConstant() {
	structures.RegisterTable("consts", &constProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	})

	structures.RegisterField("consts", "name", genericName{})
	structures.RegisterField("consts", "pkg_name", genericPackageName{})
	structures.RegisterField("consts", "pkg_path", genericPackagePath{})
	structures.RegisterField("consts", "file", genericFileName{})
	structures.RegisterField("consts", "exported", genericIsExported{})
	structures.RegisterField("consts", "docs", genericDoc{})
}

func init() {
	registerConstant()
}
