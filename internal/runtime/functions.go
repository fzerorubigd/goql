package runtime

import (
	"sync"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/fzerorubigd/goql/structures"
)

var (
	funcCache = make(map[*astdata.Package][]interface{})
	funcLock  = &sync.Mutex{}
)

func functionsProvider(in interface{}) []interface{} {
	p := in.(*astdata.Package)
	if d, ok := funcCache[p]; ok {
		return d
	}
	fn := p.Functions()
	res := make([]interface{}, len(fn))
	for i := range fn {
		res[i] = fn[i]
	}
	funcCache[p] = res
	return res
}

type isMethodColumn struct{}

func (isMethodColumn) Value(in interface{}) bool {
	fn := in.(*astdata.Function)
	return fn.Receiver() != nil
}

func registerFunc() {
	structures.RegisterTable("funcs", functionsProvider)

	structures.RegisterField("funcs", "name", genericName{})
	structures.RegisterField("funcs", "pkg_name", genericPackageName{})
	structures.RegisterField("funcs", "pkg_path", genericPackagePath{})
	structures.RegisterField("funcs", "file", genericFileName{})
	structures.RegisterField("funcs", "method", isMethodColumn{})
}

func init() {
	registerFunc()
}
