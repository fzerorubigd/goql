package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVar(t *testing.T) {
	p, err := ParsePackage("github.com/fzerorubigd/fixture")
	assert.NoError(t, err)

	f, err := p.FindVariable("hi")
	assert.NoError(t, err)

	assert.Equal(t, "hi", f.Name())
	assert.Equal(t, "// all\n// hi", f.Docs().String())
	assert.Equal(t, p, f.Package())
	assert.Equal(t, "main.go", f.File().FileName())
}
