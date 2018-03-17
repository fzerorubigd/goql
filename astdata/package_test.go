package astdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	p, err := ParsePackage("github.com/fzerorubigd/fixture")
	assert.NoError(t, err)
	assert.Equal(t, p.Path(), "github.com/fzerorubigd/fixture")
	assert.Equal(t, p.Name(), "fixture")

	_, err = translateToFullPath("invalid_path")
	assert.Error(t, err)
	_, err = translateToFullPath("github.com/fzerorubigd/package.go")
	assert.Error(t, err)
}
