package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testMap = `
package example

var aMap map[string]map[string]int
`

func TestMapType(t *testing.T) {
	p := &Package{}

	f, err := ParseFile(testMap, p)
	require.NoError(t, err)

	p.files = append(p.files, f)

	v1, err := p.FindVariable("aMap")
	require.NoError(t, err)

	require.IsType(t, &MapType{}, v1.Definition())
	mm := v1.def.(*MapType)
	assert.IsType(t, &IdentType{}, mm.Key())
	assert.Equal(t, "string", mm.key.String())
	assert.IsType(t, &MapType{}, mm.Val())
	assert.Equal(t, "map[string]int", mm.value.String())
	nd, err := NewDefinition(mm.String())
	require.NoError(t, err)
	assert.True(t, nd.Compare(mm))

	assert.Equal(t, p, mm.Package())
	assert.Equal(t, f, mm.File())
}
