package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testFunc = `
package example

func YY (in string, c bool) (int, error) {

}

func YY2 () (err error) {

}

`

func TestFuncType(t *testing.T) {
	p := &Package{}

	f, err := ParseFile(testFunc, p)
	assert.NoError(t, err)

	p.files = append(p.files, f)

	fn, err := p.FindFunction("YY")
	assert.NoError(t, err)

	def := fn.Func()
	assert.Equal(t, p, def.Package())
	assert.Equal(t, f, def.File())
	assert.Equal(t, 2, len(def.Parameters()))
	assert.Equal(t, 2, len(def.Results()))

	v1 := def.parameters[0]
	assert.Equal(t, "in", v1.Name())
	assert.IsType(t, &IdentType{}, v1.Definition())

	v2 := def.parameters[1]
	assert.Equal(t, "c", v2.Name())
	assert.IsType(t, &IdentType{}, v2.Definition())

	v3 := def.results[0]
	assert.Equal(t, "", v3.Name())
	assert.IsType(t, &IdentType{}, v3.Definition())

	v4 := def.results[1]
	assert.Equal(t, "", v4.Name())
	assert.IsType(t, &IdentType{}, v4.Definition())

	assert.Equal(t, "func (string, bool) (int, error)", def.String())

	fn, err = p.FindFunction("YY2")
	assert.NoError(t, err)
	assert.Equal(t, "func () error", fn.def.String())
}
