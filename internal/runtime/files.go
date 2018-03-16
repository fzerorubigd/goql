package runtime

import (
	"sync"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/fzerorubigd/goql/structures"
)

type filesProvider struct {
	cache map[string][]interface{}
	lock  *sync.Mutex
}

func (f *filesProvider) Provide(in interface{}) []interface{} {
	f.lock.Lock()
	defer f.lock.Unlock()

	p := in.(*astdata.Package)
	if d, ok := f.cache[p.Path()]; ok {
		return d
	}
	fs := p.Files()
	res := make([]interface{}, len(fs))
	for i := range fs {
		res[i] = fs[i]
	}
	f.cache[p.Path()] = res
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
	structures.RegisterTable("files", &filesProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	})
	structures.RegisterField("files", "name", nameColumn{})
	structures.RegisterField("files", "pkg_name", genericPackageName{})
	structures.RegisterField("files", "pkg_path", genericPackagePath{})
	structures.RegisterField("files", "docs", genericDoc{})

}

func init() {
	registerFiles()
}
