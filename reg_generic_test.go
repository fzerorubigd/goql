package goql

import (
	"sync"
	"testing"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/stretchr/testify/assert"
)

type aller struct {
	name string
	docs astdata.Docs
	pkg  *astdata.Package
	fl   *astdata.File
}

func (a aller) Name() string {
	return a.name
}

func (a aller) Docs() astdata.Docs {
	return a.docs
}

func (a aller) Package() *astdata.Package {
	return a.pkg
}

func (a aller) File() *astdata.File {
	return a.fl
}

func TestGeneric(t *testing.T) {
	assert.Equal(t, String{String: "Test"}, genericName{}.Value(aller{name: "Test"}))
	assert.Equal(t, Bool{Bool: true}, genericIsExported{}.Value(aller{name: "Test"}))
	assert.Equal(t, Bool{Bool: false}, genericIsExported{}.Value(aller{name: "test"}))
	assert.Equal(t, String{String: "// Test\n// Line2"}, genericDoc{}.Value(aller{docs: astdata.Docs{"// Test", "// Line2"}}))
	assert.Equal(t, String{Null: true}, genericDoc{}.Value(aller{}))

	p, err := astdata.ParsePackage("github.com/fzerorubigd/fixture")
	assert.NoError(t, err)
	assert.Equal(t, String{String: "fixture"}, genericPackageName{}.Value(aller{pkg: p}))
	assert.Equal(t, String{String: "github.com/fzerorubigd/fixture"}, genericPackagePath{}.Value(aller{pkg: p}))

	f := p.Files()[0]
	flName := f.FileName()
	assert.Equal(t, String{String: flName}, genericFileName{}.Value(aller{fl: f}))
	assert.Equal(t, String{String: flName}, nameColumn{}.Value(f))

	fn, err := p.FindFunction("test")
	assert.NoError(t, err)
	assert.Equal(t, String{Null: true}, isMethodColumn{}.Value(fn))
	assert.Equal(t, Bool{Null: true}, isPointerMethod{}.Value(fn))

	m1, err := p.FindMethod("beta", "assert")
	assert.NoError(t, err)
	assert.Equal(t, String{String: "beta"}, isMethodColumn{}.Value(m1))
	assert.Equal(t, Bool{Bool: false}, isPointerMethod{}.Value(m1))

	m2, err := p.FindMethod("alpha", "testing")
	assert.NoError(t, err)
	assert.Equal(t, String{String: "alpha"}, isMethodColumn{}.Value(m2))
	assert.Equal(t, Bool{Bool: true}, isPointerMethod{}.Value(m2))
	assert.Equal(t, String{String: "\n\tpanic(\"hi\")\n"}, bodyCol{}.Value(m2))

	im1, err := p.FindImport("net/http")
	assert.NoError(t, err)
	assert.Equal(t, String{Null: true}, canonicalCol{}.Value(im1))
	assert.Equal(t, String{String: "net/http"}, pathCol{}.Value(im1))
	assert.Equal(t, String{String: "http"}, packageCol{}.Value(im1))

	im2, err := p.FindImport("context")
	assert.NoError(t, err)
	assert.Equal(t, String{String: "ctx"}, canonicalCol{}.Value(im2))
	assert.Equal(t, String{String: "context"}, pathCol{}.Value(im2))
	assert.Equal(t, String{String: "context"}, packageCol{}.Value(im2))

}

func TestProviders(t *testing.T) {
	p, err := astdata.ParsePackage("github.com/fzerorubigd/fixture")
	assert.NoError(t, err)

	fn := &functionProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	}

	data := fn.Provide(p)
	for i := range data {
		assert.IsType(t, &astdata.Function{}, data[i])
	}
	// cache?
	data = fn.Provide(p)
	for i := range data {
		assert.IsType(t, &astdata.Function{}, data[i])
	}

	vn := &typeProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	}

	data = vn.Provide(p)
	for i := range data {
		assert.IsType(t, &astdata.Type{}, data[i])
	}
	// cache?
	data = vn.Provide(p)
	for i := range data {
		assert.IsType(t, &astdata.Type{}, data[i])
	}

	vn2 := &variableProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	}

	data = vn2.Provide(p)
	for i := range data {
		assert.IsType(t, &astdata.Variable{}, data[i])
	}
	// cache?
	data = vn2.Provide(p)
	for i := range data {
		assert.IsType(t, &astdata.Variable{}, data[i])
	}

	vn3 := &filesProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	}

	data = vn3.Provide(p)
	for i := range data {
		assert.IsType(t, &astdata.File{}, data[i])
	}
	// cache?
	data = vn3.Provide(p)
	for i := range data {
		assert.IsType(t, &astdata.File{}, data[i])
	}

	vn4 := &constProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	}

	data = vn4.Provide(p)
	for i := range data {
		assert.IsType(t, &astdata.Constant{}, data[i])
	}
	// cache?
	data = vn4.Provide(p)
	for i := range data {
		assert.IsType(t, &astdata.Constant{}, data[i])
	}

	vn5 := &importProvider{
		cache: make(map[string][]interface{}),
		lock:  &sync.Mutex{},
	}

	data = vn5.Provide(p)
	for i := range data {
		assert.IsType(t, &astdata.Import{}, data[i])
	}
	// cache?
	data = vn5.Provide(p)
	for i := range data {
		assert.IsType(t, &astdata.Import{}, data[i])
	}
}
