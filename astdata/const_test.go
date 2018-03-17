package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConst(t *testing.T) {
	p, err := ParsePackage("github.com/fzerorubigd/fixture")
	assert.NoError(t, err)

	c, err := p.FindConstant("X")
	assert.NoError(t, err)

	assert.Equal(t, "X", c.Name())
	assert.Equal(t, "// X Docs", c.Docs().String())
	assert.Equal(t, p, c.Package())
	assert.Equal(t, "main.go", c.File().FileName())

	c, err = p.FindConstant("testConst")
	assert.NoError(t, err)

	assert.Equal(t, "testConst", c.Name())
	assert.Equal(t, "", c.Docs().String())
	assert.Equal(t, p, c.Package())
	assert.Equal(t, "main.go", c.File().FileName())
	assert.Equal(t, "10", c.Value())

}
