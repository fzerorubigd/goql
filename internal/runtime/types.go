package runtime

import (
	"sync"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/fzerorubigd/goql/structures"
)

var (
	typesCache = make(map[*astdata.Package][]interface{})
	typesLock  = &sync.Mutex{}
)

func typesProvider(in interface{}) []interface{} {
	typesLock.Lock()
	defer typesLock.Unlock()

	p := in.(*astdata.Package)
	if d, ok := typesCache[p]; ok {
		return d
	}
	fs := p.Types()
	res := make([]interface{}, len(fs))
	for i := range fs {
		res[i] = fs[i]
	}
	typesCache[p] = res
	return res
}

func registerTypes() {
	structures.RegisterTable("types", typesProvider)

	structures.RegisterField("types", "name", genericName{})
	structures.RegisterField("types", "pkg_name", genericPackageName{})
	structures.RegisterField("types", "pkg_path", genericPackagePath{})
	structures.RegisterField("types", "file", genericFileName{})
	structures.RegisterField("types", "exported", genericIsExported{})
}

func init() {
	registerTypes()
}
