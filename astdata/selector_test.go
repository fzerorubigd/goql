package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testSel = `
package example

import (
"context"
"net/http"
)

type ZZ struct {
    C context.Context
    R *http.Request
}
`

func TestSelectorType(t *testing.T) {
	p := &Package{}

	f, err := ParseFile(testSel, p)
	assert.NoError(t, err)

	p.files = append(p.files, f)

	tt, err := p.FindType("ZZ")
	assert.NoError(t, err)

	assert.IsType(t, &StructType{}, tt.def)
	st := tt.def.(*StructType)

	assert.Equal(t, 2, len(st.fields))
	v1 := st.fields[0].def
	assert.IsType(t, &SelectorType{}, v1)
	sel := v1.(*SelectorType)
	assert.Equal(t, "context", sel.Selector())

	imp, err := f.FindImport("context")
	assert.NoError(t, err)
	assert.Equal(t, imp, sel.Import())
	assert.Equal(t, p, sel.Package())
	assert.Equal(t, "Context", sel.Ident())
	assert.Equal(t, "context.Context", sel.String())
	nd, err := NewDefinition(v1.String())
	require.NoError(t, err)
	assert.True(t, nd.Compare(v1))

}
