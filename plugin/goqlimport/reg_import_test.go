package goqlimport

import (
	"testing"

	"github.com/fzerorubigd/goql/astdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fl1 = `package example
`

var fl2 = `package example

var (
notFormated = true
)
`

func TestImport(t *testing.T) {
	ic := newImportCheck()
	p := &astdata.Package{}
	fl1F, err := astdata.ParseFile(fl1, p)
	require.NoError(t, err)
	assert.True(t, ic.Value(fl1F).Bool)

	fl2F, err := astdata.ParseFile(fl2, p)
	require.NoError(t, err)
	assert.False(t, ic.Value(fl2F).Bool)

	// Try to register it
	assert.NotPanics(t, Register)
	assert.Panics(t, Register)
}
