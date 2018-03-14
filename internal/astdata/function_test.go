package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunction(t *testing.T) {
	p, err := ParsePackage("github.com/fzerorubigd/goql/internal/astdata/fixture")
	assert.NoError(t, err)

	f, err := p.FindFunction("test")
	assert.NoError(t, err)

	assert.Equal(t, "test", f.Name())
	assert.Equal(t, "// Multi line\n// comment\n// here", f.Docs().String())
	assert.Equal(t, p, f.Package())
	assert.Equal(t, "main.go", f.File().FileName())

	f, err = p.FindFunction("beta.assert")
	assert.NoError(t, err)

	assert.Equal(t, "beta.assert", f.Name())
	assert.Equal(t, "", f.Docs().String())
	assert.Equal(t, p, f.Package())
	assert.Equal(t, "main.go", f.File().FileName())
	assert.NotNil(t, f.Receiver())
	assert.Equal(t, "", f.Receiver().Name())

	f, err = p.FindFunction("alpha.testing")
	assert.NoError(t, err)

	assert.Equal(t, "alpha.testing", f.Name())
	assert.Equal(t, "", f.Docs().String())
	assert.Equal(t, p, f.Package())
	assert.Equal(t, "main.go", f.File().FileName())
	assert.NotNil(t, f.Receiver())
	assert.Equal(t, "a", f.Receiver().Name())

}
