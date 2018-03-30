package goql

import (
	"sync"

	"github.com/fzerorubigd/goql/astdata"
)

type functionProvider struct {
	cache map[string][]interface{}
	lock  *sync.Mutex
}

func (f *functionProvider) Provide(in interface{}) []interface{} {
	f.lock.Lock()
	defer f.lock.Unlock()

	p := in.(*astdata.Package)
	if d, ok := f.cache[p.Path()]; ok {
		return d
	}
	fn := p.Functions()
	res := make([]interface{}, len(fn))
	for i := range fn {
		res[i] = fn[i]
	}
	f.cache[p.Path()] = res
	return res
}

type isMethodColumn struct{}

func (isMethodColumn) Value(in interface{}) String {
	fn := in.(*astdata.Function)
	if fn.ReceiverType() == "" {
		return String{Null: true}
	}
	return String{String: fn.ReceiverType()}
}

type isPointerMethod struct{}

func (isPointerMethod) Value(in interface{}) Bool {
	fn := in.(*astdata.Function)
	if fn.ReceiverType() == "" {
		return Bool{Null: true}
	}
	return Bool{Bool: fn.RecieverPointer()}
}

type bodyCol struct{}

func (bodyCol) Value(in interface{}) String {
	fn := in.(*astdata.Function)
	return String{String: fn.Body()}
}

func registerFunc() {
	RegisterTable("funcs", &functionProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	})

	RegisterField("funcs", "name", genericName{})
	RegisterField("funcs", "pkg_name", genericPackageName{})
	RegisterField("funcs", "pkg_path", genericPackagePath{})
	RegisterField("funcs", "file", genericFileName{})
	RegisterField("funcs", "receiver", isMethodColumn{})
	RegisterField("funcs", "pointer_receiver", isPointerMethod{})
	RegisterField("funcs", "exported", genericIsExported{})
	RegisterField("funcs", "docs", genericDoc{})
	RegisterField("funcs", "body", bodyCol{})
}

func init() {
	registerFunc()
}
