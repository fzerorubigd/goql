package runtime

import (
	"sync"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/fzerorubigd/goql/structures"
)

var (
	filesCache = make(map[*astdata.Package][]interface{})
	fileLock   = &sync.Mutex{}
)

func filesProvider(in interface{}) []interface{} {
	fileLock.Lock()
	defer fileLock.Unlock()

	p := in.(*astdata.Package)
	if d, ok := filesCache[p]; ok {
		return d
	}
	fs := p.Files()
	res := make([]interface{}, len(fs))
	for i := range fs {
		res[i] = fs[i]
	}
	filesCache[p] = res
	return res
}

type nameColumn struct {
}

func (nameColumn) Value(in interface{}) structures.String {
	fl := in.(*astdata.File)
	return structures.String{String: fl.FileName()}
}

func registerFiles() {
	// register files
	structures.RegisterTable("files", filesProvider)
	structures.RegisterField("files", "name", nameColumn{})
	structures.RegisterField("files", "pkg_name", genericPackageName{})
	structures.RegisterField("files", "pkg_path", genericPackagePath{})
	structures.RegisterField("files", "docs", genericDoc{})

}

func init() {
	registerFiles()
}
