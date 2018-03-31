package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testVar = `
package example

import "context"

var (
    vi = 10
    vj = 100.0
    vk = 100i
    vl = "str"
    vm = 'c'
    vn = Test()
    vo = make([]int,10)
    vp = new([]int)
    vq = context.Background()
)

func Test() int64 {
    return 10
}
`

func TestVar(t *testing.T) {
	p, err := ParsePackage("github.com/fzerorubigd/fixture")
	require.NoError(t, err)

	f, err := p.FindVariable("hi")
	assert.NoError(t, err)

	assert.Equal(t, "hi", f.Name())
	assert.Equal(t, "// all\n// hi", f.Docs().String())
	assert.Equal(t, p, f.Package())
	assert.Equal(t, "main.go", f.File().FileName())
}

func TestVarExtra(t *testing.T) {
	p := &Package{}
	f, err := ParseFile(testVar, p)
	require.NoError(t, err)
	p.files = append(p.files, f)

	v, err := p.FindVariable("vo")
	require.NoError(t, err)

	assert.Equal(t, "[]int", v.Definition().String())

	v, err = p.FindVariable("vn")
	require.NoError(t, err)

	assert.Equal(t, "int64", v.Definition().String())

	v, err = p.FindVariable("vp")
	require.NoError(t, err)

	assert.Equal(t, "*[]int", v.Definition().String())

	v, err = p.FindVariable("vq")
	require.NoError(t, err)

	assert.Equal(t, "context.Context", v.Definition().String())

}
