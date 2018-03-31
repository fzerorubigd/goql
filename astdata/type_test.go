package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestType(t *testing.T) {
	p, err := ParsePackage("github.com/fzerorubigd/fixture")
	require.NoError(t, err)

	f, err := p.FindType("alpha")
	assert.NoError(t, err)

	assert.Equal(t, "alpha", f.Name())
	assert.Equal(t, "// aaa\n// alpha comment", f.Docs().String())
	assert.Equal(t, p, f.Package())
	assert.Equal(t, "main.go", f.File().FileName())

}
