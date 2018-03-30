package goql

import (
	"sync"

	"github.com/fzerorubigd/goql/astdata"
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

func (nameColumn) Value(in interface{}) String {
	fl := in.(*astdata.File)
	return String{String: fl.FileName()}
}

func registerFiles() {
	// register files
	RegisterTable("files", &filesProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	})
	RegisterField("files", "name", nameColumn{})
	RegisterField("files", "pkg_name", genericPackageName{})
	RegisterField("files", "pkg_path", genericPackagePath{})
	RegisterField("files", "docs", genericDoc{})

}

func init() {
	registerFiles()
}
