package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testConst = `
package example

import "math"

const (
   ci = 10
   cj = 100.0
   ck = 10i
   cl string = "str"
   cm = 'c'
   cn = "another"

   pi = math.Pi
)

const (
   eci = iota
   ecj
)

`

func TestConst(t *testing.T) {
	p, err := ParsePackage("github.com/fzerorubigd/fixture")
	assert.NoError(t, err)

	c, err := p.FindConstant("X")
	assert.NoError(t, err)

	assert.Equal(t, "X", c.Name())
	assert.Equal(t, "// X Docs", c.Docs().String())
	assert.Equal(t, p, c.Package())
	//	assert.Equal(t, "main.go", c.File().FileName())

	c, err = p.FindConstant("testConst")
	assert.NoError(t, err)

	assert.Equal(t, "testConst", c.Name())
	assert.Equal(t, "", c.Docs().String())
	assert.Equal(t, p, c.Package())
	//	assert.Equal(t, "main.go", c.File().FileName())
	assert.Equal(t, "10", c.Value())

	p = &Package{}
	f, err := ParseFile(testConst, p)
	require.NoError(t, err)

	p.files = append(p.files, f)

	c, err = p.FindConstant("ci")
	require.NoError(t, err)
	assert.Equal(t, f, c.File())
	assert.Equal(t, p, c.Package())
}
