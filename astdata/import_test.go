package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	p, err := ParsePackage("github.com/fzerorubigd/fixture")
	assert.NoError(t, err)

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
}
