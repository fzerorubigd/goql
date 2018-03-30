package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testArrr = `
package example

type A []string

var X = [...]string{"a","b","c"}

var Y = [10]string{}

`

func TestArrayType(t *testing.T) {
	p := &Package{}

	f, err := ParseFile(testArrr, p)
	assert.NoError(t, err)

	p.files = append(p.files, f)
	typ, err := p.FindType("A")
	assert.NoError(t, err)
	assert.Equal(t, "A", typ.Name())

	def := typ.Definition()
	assert.IsType(t, &ArrayType{}, def)
	at := def.(*ArrayType)
	assert.True(t, at.Slice())
	assert.Equal(t, p, at.Package())
	assert.Equal(t, f, at.File())
	assert.Equal(t, 0, at.Len())
	assert.IsType(t, &IdentType{}, at.ValueDefinition())

	assert.Equal(t, "[]string", at.String())

	v, err := p.FindVariable("X")
	assert.NoError(t, err)
	def = v.Definition()
	assert.IsType(t, &EllipsisType{}, def)
	et := def.(*EllipsisType)
	assert.False(t, et.Slice())
	assert.Equal(t, p, et.Package())
	assert.Equal(t, f, et.File())
	assert.Equal(t, 0, et.Len()) // TODO : size
	assert.IsType(t, &IdentType{}, et.ValueDefinition())

	assert.Equal(t, "[...]string", et.String())

	v, err = p.FindVariable("Y")
	assert.NoError(t, err)
	def = v.Definition()
	assert.IsType(t, &ArrayType{}, def)
	at = def.(*ArrayType)
	assert.False(t, at.Slice())
	assert.Equal(t, p, at.Package())
	assert.Equal(t, f, at.File())
	assert.Equal(t, 10, at.Len())
	assert.IsType(t, &IdentType{}, at.ValueDefinition())

	assert.Equal(t, "[10]string", at.String())

}
