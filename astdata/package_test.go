package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackage(t *testing.T) {
	p, err := ParsePackage("github.com/fzerorubigd/fixture")
	assert.NoError(t, err)
	assert.Equal(t, p.Path(), "github.com/fzerorubigd/fixture")
	assert.Equal(t, p.Name(), "fixture")

	_, err = translateToFullPath("invalid_path")
	assert.Error(t, err)
	_, err = translateToFullPath("github.com/fzerorubigd/not/exists/package")
	assert.Error(t, err)

	p = &Package{}

	f1, err := ParseFile(testArrr, p)
	require.NoError(t, err)
	f2, err := ParseFile(testFunc, p)
	require.NoError(t, err)
	f3, err := ParseFile(testImport, p)
	require.NoError(t, err)
	f4, err := ParseFile(testConst, p)
	require.NoError(t, err)

	p.files = append(p.files, f1, f2, f3, f4)

	assert.Len(t, p.Files(), 4)
	assert.NotEmpty(t, p.Functions())
	assert.NotEmpty(t, p.Variables())
	assert.NotEmpty(t, p.Types())
	assert.NotEmpty(t, p.Imports())
	assert.NotEmpty(t, p.Constants())
}
