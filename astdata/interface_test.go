package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testInterface = `
package example

type II interface {
    Test(string)error
}

type II2 interface {
    II
    Hi(string)string
}

type II3 interface{}
`

func TestInterfaceType(t *testing.T) {
	p := &Package{}
	f, err := ParseFile(testInterface, p)
	assert.NoError(t, err)

	p.files = append(p.files, f)

	tn, err := p.FindType("II")
	assert.NoError(t, err)

	require.IsType(t, tn.def, &InterfaceType{})
	def := tn.def.(*InterfaceType)

	require.Len(t, def.Functions(), 1)
	require.Len(t, def.Embeds(), 0)

	fn := def.Functions()[0]
	assert.Equal(t, "Test", fn.Name())
	assert.Equal(t, "func (string) error", fn.Func().String())
	assert.Equal(t, p, def.Package())
	assert.Equal(t, f, def.File())

	assert.Equal(t, "interface{\n\tTest(string) error\n}", def.String())

	tn, err = p.FindType("II2")
	assert.NoError(t, err)

	require.IsType(t, tn.def, &InterfaceType{})
	def = tn.def.(*InterfaceType)

	require.Len(t, def.Functions(), 1)
	require.Len(t, def.Embeds(), 1)

	fn = def.Functions()[0]
	assert.Equal(t, "Hi", fn.Name())
	assert.Equal(t, "func (string) string", fn.Func().String())
	assert.Equal(t, p, def.Package())
	assert.Equal(t, f, def.File())

	em := def.Embeds()[0]
	require.IsType(t, &IdentType{}, em)
	assert.Equal(t, "II", em.String())

	assert.Equal(t, "interface{\n\tII\n\tHi(string) string\n}", def.String())

	tn, err = p.FindType("II3")
	assert.NoError(t, err)

	require.IsType(t, tn.def, &InterfaceType{})
	def = tn.def.(*InterfaceType)

	require.Len(t, def.Functions(), 0)
	require.Len(t, def.Embeds(), 0)

	assert.Equal(t, "interface{}", def.String())

}
