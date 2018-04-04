package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testStar = `
package example

var intStar *int
`

func TestStarType(t *testing.T) {
	p := &Package{}

	f, err := ParseFile(testStar, p)
	assert.NoError(t, err)
	p.files = append(p.files, f)

	tn, err := p.FindVariable("intStar")
	assert.NoError(t, err)

	assert.IsType(t, &StarType{}, tn.Definition())
	s := tn.def.(*StarType)
	assert.IsType(t, &IdentType{}, s.Target())
	assert.Equal(t, "*int", s.String())
	nd, err := NewDefinition(s.String())
	require.NoError(t, err)
	assert.True(t, nd.Compare(s))

	assert.Equal(t, p, s.Package())
	assert.Equal(t, f, s.File())
}
