package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testIdent = `
package example

type ALPHA string

var x ALPHA

`

func TestIdentType(t *testing.T) {
	p := &Package{}

	f, err := ParseFile(testIdent, p)
	assert.NoError(t, err)
	p.files = append(p.files, f)

	tp, err := p.FindType("ALPHA")
	assert.NoError(t, err)

	assert.IsType(t, &IdentType{}, tp.def)
	def := tp.def.(*IdentType)
	assert.Equal(t, "string", def.Ident())
	assert.Equal(t, "string", def.String())
	assert.Equal(t, p, def.Package())
	assert.Equal(t, f, def.File())
	nd, err := NewDefinition(def.String())
	require.NoError(t, err)
	assert.True(t, nd.Compare(def))

	v, err := p.FindVariable("x")
	assert.NoError(t, err)

	assert.IsType(t, &IdentType{}, v.def)
	def = v.def.(*IdentType)
	assert.Equal(t, "ALPHA", def.Ident())
	assert.Equal(t, "ALPHA", def.String())
	nd, err = NewDefinition(def.String())
	require.NoError(t, err)
	assert.True(t, nd.Compare(def))

}
