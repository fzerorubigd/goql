package goql

import (
	"sync"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/fzerorubigd/goql/structures"
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

func (isMethodColumn) Value(in interface{}) structures.String {
	fn := in.(*astdata.Function)
	if fn.ReceiverType() == "" {
		return structures.String{Null: true}
	}
	return structures.String{String: fn.ReceiverType()}
}

type isPointerMethod struct{}

func (isPointerMethod) Value(in interface{}) structures.Bool {
	fn := in.(*astdata.Function)
	if fn.ReceiverType() == "" {
		return structures.Bool{Null: true}
	}
	return structures.Bool{Bool: fn.RecieverPointer()}
}

type bodyCol struct{}

func (bodyCol) Value(in interface{}) structures.String {
	fn := in.(*astdata.Function)
	return structures.String{String: fn.Body()}
}

func registerFunc() {
	structures.RegisterTable("funcs", &functionProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	})

	structures.RegisterField("funcs", "name", genericName{})
	structures.RegisterField("funcs", "pkg_name", genericPackageName{})
	structures.RegisterField("funcs", "pkg_path", genericPackagePath{})
	structures.RegisterField("funcs", "file", genericFileName{})
	structures.RegisterField("funcs", "receiver", isMethodColumn{})
	structures.RegisterField("funcs", "pointer_receiver", isPointerMethod{})
	structures.RegisterField("funcs", "exported", genericIsExported{})
	structures.RegisterField("funcs", "docs", genericDoc{})
	structures.RegisterField("funcs", "body", bodyCol{})
}

func init() {
	registerFunc()
}
