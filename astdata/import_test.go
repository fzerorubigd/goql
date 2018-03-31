package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testImport = `
package example

import "github.com/fzerorubigd/fixture"
`

func TestImport(t *testing.T) {
	p, err := ParsePackage("github.com/fzerorubigd/fixture")
	require.NoError(t, err)

	f, err := p.FindImport("ctx")
	assert.NoError(t, err)

	assert.Equal(t, "ctx", f.Canonical())
	assert.Equal(t, "context", f.TargetPackage())
	assert.Regexp(t, "context$", f.Folder())
	assert.Equal(t, "// comment ctx", f.Docs().String())
	assert.Equal(t, p, f.Package())
	assert.Equal(t, "main.go", f.File().FileName())

	f, err = p.FindImport("net/http")
	assert.NoError(t, err)

	assert.Equal(t, "", f.Canonical())
	assert.Equal(t, "http", f.TargetPackage())
	assert.Regexp(t, "net/http$", f.Folder())
	assert.Equal(t, "// comment http", f.Docs().String())
	assert.Equal(t, p, f.Package())
	assert.Equal(t, "main.go", f.File().FileName())

	p = &Package{}
	fl, err := ParseFile(testImport, p)
	require.NoError(t, err)

	p.files = append(p.files, fl)

	for i := range fl.imports {
		assert.Regexp(t, "github.com/fzerorubigd/fixture$", fl.imports[i].Folder())
	}
}
