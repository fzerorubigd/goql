package runtime

import (
	"sync"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/fzerorubigd/goql/structures"
)

var (
	varCache = make(map[*astdata.Package][]interface{})
	varLock  = &sync.Mutex{}
)

func variablesProvider(in interface{}) []interface{} {
	p := in.(*astdata.Package)
	if d, ok := varCache[p]; ok {
		return d
	}
	v := p.Variables()
	res := make([]interface{}, len(v))
	for i := range v {
		res[i] = v[i]
	}
	varCache[p] = res
	return res
}

func registerVariable() {
	structures.RegisterTable("vars", variablesProvider)

	structures.RegisterField("vars", "name", genericName{})
	structures.RegisterField("vars", "pkg_name", genericPackageName{})
	structures.RegisterField("vars", "pkg_path", genericPackagePath{})
	structures.RegisterField("vars", "file", genericFileName{})
	structures.RegisterField("vars", "exported", genericIsExported{})
}

func init() {
	registerVariable()
}
