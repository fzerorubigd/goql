package goql

import (
	"sync"

	"github.com/fzerorubigd/goql/astdata"
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
	RegisterTable("consts", &constProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	})

	RegisterField("consts", "name", genericName{})
	RegisterField("consts", "pkg_name", genericPackageName{})
	RegisterField("consts", "pkg_path", genericPackagePath{})
	RegisterField("consts", "file", genericFileName{})
	RegisterField("consts", "exported", genericIsExported{})
	RegisterField("consts", "docs", genericDoc{})
}

func init() {
	registerConstant()
}
